package api

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
)

func (a *App) CreateKitchenRecipe(c echo.Context) error {
	var input recipes.Recipe
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)
	kitchenID := c.Param("kitchen_id")

	// Validate that the user has permission to this kitchen.
	// Lookup the user record for the provided JWT.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account by user ID").SetInternal(err)
	}

	kitchen, err := kitchens.GetKitchenByID(ctx, a.db, kitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get kitchen by ID").SetInternal(err)
	}

	// Validate that the user has permissions to be modifying this kitchen.
	if account.AccountID != kitchen.AccountID {
		err := fmt.Errorf("account '%s' attempted to create recipe in '%s' without authorization", account.AccountID, kitchen.KitchenID)
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}

	// Create recipe ID
	recipeID := recipes.CreateRecipeID()

	// Handle the file uploads if they have been set.
	var recipe recipes.Recipe

	err = mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		recipe, err = recipes.CreateRecipe(ctx, tx, recipes.CreateRecipeInput{
			RecipeID:  recipeID,
			KitchenID: kitchenID,
			Name:      input.Name,
			Summary:   input.Summary,
			PrepTime:  *input.PrepTime,
			CookTime:  *input.CookTime,
			Servings:  *input.Servings,
			Cover:     input.Cover,
			Source:    input.Source,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err).SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipe)
}

func (a *App) GetKitchenRecipe(c echo.Context) error {
	ctx := c.Request().Context()
	recipeID := c.Param("recipe_id")

	recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe by ID").SetInternal(err)
	}

	steps, err := recipes.GetRecipeStepsByRecipeID(ctx, a.db, recipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe steps").SetInternal(err)
	}

	// If there are no steps, return early.
	if len(steps) == 0 {
		return c.JSON(http.StatusOK, recipe)
	}

	recipe.Steps = steps

	// TODO: Add images and notes to steps.

	return c.JSON(http.StatusOK, recipe)
}

func (a *App) GetKitchenRecipes(c echo.Context) error {
	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	recipes, err := recipes.ListRecipesByKitchenID(ctx, a.db, kitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe(s) by kitchen ID").SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipes)
}
