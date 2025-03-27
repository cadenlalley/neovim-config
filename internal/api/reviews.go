package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type CreateRecipeReviewRequest struct {
	Description string   `json:"description"`
	Rating      *float64 `json:"rating" validate:"required"`
	Media       string   `json:"media"`
}

func (a *App) CreateRecipeReview(c echo.Context) error {
	var input CreateRecipeReviewRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate the input
	if len(input.Description) > 2000 {
		return echo.NewHTTPError(http.StatusBadRequest, "description must be less than 2000 characters")
	}

	if input.Rating == nil || *input.Rating < 1 || *input.Rating > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "rating must be between 1 and 5")
	}

	ctx := c.Request().Context()
	userId := c.Get(auth.UserIDContextKey).(string)
	userKitchenID := c.Get(auth.KitchenIDContextKey).(string)
	recipeID := c.Param("recipe_id")

	// Check that the user is allowed to create a review for the input kitchen.
	allowed, err := kitchens.CheckKitchenWriter(ctx, a.db, userKitchenID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify user permissions to kitchen").SetInternal(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	reviewID := recipes.CreateReviewID()
	review, err := recipes.CreateReview(ctx, a.db, recipes.CreateReviewInput{
		ReviewID:          reviewID,
		RecipeID:          recipeID,
		ReviewerKitchenID: userKitchenID,
		Description:       null.NewString(input.Description, input.Description != ""),
		Rating:            *input.Rating,
		Media:             null.NewString(input.Media, input.Media != ""),
	})
	if err != nil {
		if err == recipes.ErrDuplicateReview {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create recipe review").SetInternal(err)
	}

	return c.JSON(http.StatusOK, review)
}

type UpdateRecipeReviewRequest struct {
	Description *string  `json:"description"`
	Rating      *float64 `json:"rating" validate:"required"`
	Media       *string  `json:"media"`
}

func (a *App) UpdateRecipeReview(c echo.Context) error {
	var input UpdateRecipeReviewRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate the input
	if input.Description != nil && len(*input.Description) > 2000 {
		return echo.NewHTTPError(http.StatusBadRequest, "description must be less than 2000 characters")
	}

	if input.Rating == nil || *input.Rating < 1 || *input.Rating > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "rating must be between 1 and 5")
	}

	ctx := c.Request().Context()
	userId := c.Get(auth.UserIDContextKey).(string)
	reviewID := c.Param("review_id")

	// Get the review to update.
	review, err := recipes.GetReviewByID(ctx, a.db, reviewID)
	if err != nil {
		if err == recipes.ErrReviewNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe review").SetInternal(err)
	}

	// Check that the user is allowed to update a review for the input kitchen.
	allowed, err := kitchens.CheckKitchenWriter(ctx, a.db, review.ReviewerKitchenID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify user permissions to kitchen").SetInternal(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Nullable inputs
	if input.Description != nil && *input.Description != review.Description.String {
		review.Description = null.NewString(*input.Description, *input.Description != "")
	}

	if input.Rating != nil && *input.Rating != review.Rating {
		review.Rating = *input.Rating
	}

	if input.Media != nil && *input.Media != review.Media.String {
		review.Media = null.NewString(*input.Media, *input.Media != "")
	}

	review, err = recipes.UpdateReview(ctx, a.db, recipes.UpdateReviewInput{
		ReviewID:    reviewID,
		Description: review.Description,
		Rating:      review.Rating,
		Media:       review.Media,
	})
	if err != nil {
		if err == recipes.ErrReviewNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update recipe review").SetInternal(err)
	}

	return c.JSON(http.StatusOK, review)
}

func (a *App) GetRecipeReviews(c echo.Context) error {
	ctx := c.Request().Context()
	userKitchenID := c.Get(auth.KitchenIDContextKey).(string)
	recipeID := c.Param("recipe_id")

	reviewSummary, err := recipes.GetReviewSummaryByRecipeID(ctx, a.db, recipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get review summary").SetInternal(err)
	}

	reviews, err := recipes.ListReviewsByRecipeID(ctx, a.db, recipeID, userKitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get reviews for recipe").SetInternal(err)
	}

	reviewSummary.Reviews = reviews

	// Identify the review for the active kitchen.
	for _, review := range reviews {
		if review.ReviewerKitchenID == userKitchenID {
			reviewSummary.KitchenReview = review
		}
	}

	return c.JSON(http.StatusOK, reviewSummary)
}

func (a *App) DeleteRecipeReview(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(auth.UserIDContextKey).(string)
	reviewId := c.Param("review_id")

	// Get the review to delete.
	review, err := recipes.GetReviewByID(ctx, a.db, reviewId)
	if err != nil {
		if err == recipes.ErrReviewNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get review").SetInternal(err)
	}

	// Check that the user is allowed to delete the review.
	allowed, err := kitchens.CheckKitchenWriter(ctx, a.db, review.ReviewerKitchenID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify user permissions to kitchen").SetInternal(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		err = recipes.DeleteLikesByReviewID(ctx, tx, review.ReviewID)
		if err != nil {
			return errors.Wrap(err, "could not delete recipe review likes")
		}

		err = recipes.DeleteReview(ctx, tx, review.ReviewID)
		if err != nil {
			if err == recipes.ErrReviewNotFound {
				return nil
			}
			return errors.Wrap(err, "could not delete recipe review")
		}

		return nil
	})
	if txErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete recipe review").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (a *App) LikeRecipeReview(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(auth.UserIDContextKey).(string)
	userKitchenID := c.Get(auth.KitchenIDContextKey).(string)
	reviewID := c.Param("review_id")

	// Verify access to unlike this review.
	allowed, err := kitchens.CheckKitchenWriter(ctx, a.db, userKitchenID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify user permissions to kitchen").SetInternal(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Like the review.
	err = recipes.LikeReview(ctx, a.db, reviewID, userKitchenID)
	if err != nil {
		if err == recipes.ErrReviewNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not like recipe review").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (a *App) UnlikeRecipeReview(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(auth.UserIDContextKey).(string)
	userKitchenID := c.Get(auth.KitchenIDContextKey).(string)
	reviewID := c.Param("review_id")

	// Verify access to unlike this review.
	allowed, err := kitchens.CheckKitchenWriter(ctx, a.db, userKitchenID, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not verify user permissions to kitchen").SetInternal(err)
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Unlike the review.
	err = recipes.UnlikeReview(ctx, a.db, reviewID, userKitchenID)
	if err != nil {
		if err == recipes.ErrReviewNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not dislike recipe review").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
