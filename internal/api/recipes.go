package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
)

func (a *App) CreateKitchenRecipe(c echo.Context) error {
	var input recipes.Recipe
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Additional validation that ingredients and steps were required.
	err = input.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	// Create recipe ID
	recipeID := recipes.CreateRecipeID()
	if input.RecipeID != "" {
		recipeID = input.RecipeID
	}

	var recipe recipes.Recipe
	txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
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

		// Handle step processing (images, notes)
		for _, step := range input.Steps {
			err = recipes.CreateRecipeSteps(ctx, tx, recipes.CreateRecipeStepInput{
				RecipeID:    recipeID,
				StepID:      step.StepID,
				Instruction: step.Instruction,
				Group:       step.Group,
			})
			if err != nil {
				return err
			}

			if len(step.Images) != 0 {
				for _, image := range step.Images {
					err = recipes.CreateRecipeImages(ctx, tx, recipes.CreateRecipeImagesInput{
						RecipeID: recipeID,
						StepID:   step.StepID,
						ImageURL: image,
					})
					if err != nil {
						return err
					}
				}
			}

			if step.Note != "" {
				err = recipes.CreateRecipeNotes(ctx, tx, recipes.CreateRecipeNotesInput{
					RecipeID: recipeID,
					StepID:   step.StepID,
					Note:     step.Note,
				})
				if err != nil {
					return err
				}
			}
		}

		// Handle ingredient processing
		for _, ingredient := range input.Ingredients {
			err = recipes.CreateRecipeIngredients(ctx, tx, recipes.CreateRecipeIngredientInput{
				RecipeID:     recipeID,
				IngredientID: ingredient.IngredientID,
				Name:         ingredient.Name,
				Quantity:     ingredient.Quantity,
				Unit:         ingredient.Unit,
				Group:        ingredient.Group,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if txErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create recipe").SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipe)
}

func (a *App) UpdateKitchenRecipe(c echo.Context) error {
	var input recipes.Recipe
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	recipeID := c.Param("recipe_id")

	recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		if err == recipes.ErrRecipeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe by ID").SetInternal(err)
	}

	txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		recipe, err = recipes.UpdateRecipe(ctx, tx, recipes.UpdateRecipeInput{
			RecipeID: recipeID,
			Name:     input.Name,
			Summary:  input.Summary,
			PrepTime: *input.PrepTime,
			CookTime: *input.CookTime,
			Servings: *input.Servings,
			Cover:    input.Cover,
			Source:   input.Source,
		})
		if err != nil {
			return err
		}

		// Delete all existing references for recipe steps for a full replacement.
		err = recipes.DeleteRecipeImagesByRecipeID(ctx, tx, recipeID)
		if err != nil {
			return err
		}

		err = recipes.DeleteRecipeNotesByRecipeID(ctx, tx, recipeID)
		if err != nil {
			return err
		}

		err = recipes.DeleteRecipeStepsByRecipeID(ctx, tx, recipeID)
		if err != nil {
			return err
		}

		err = recipes.DeleteRecipeIngredientsByRecipeID(ctx, tx, recipeID)
		if err != nil {
			return err
		}

		// Handle step processing (images, notes)
		for _, step := range input.Steps {
			err = recipes.CreateRecipeSteps(ctx, tx, recipes.CreateRecipeStepInput{
				RecipeID:    recipeID,
				StepID:      step.StepID,
				Instruction: step.Instruction,
				Group:       step.Group,
			})
			if err != nil {
				return err
			}

			if len(step.Images) != 0 {
				for _, image := range step.Images {
					err = recipes.CreateRecipeImages(ctx, tx, recipes.CreateRecipeImagesInput{
						RecipeID: recipeID,
						StepID:   step.StepID,
						ImageURL: image,
					})
					if err != nil {
						return err
					}
				}
			}

			if step.Note != "" {
				err = recipes.CreateRecipeNotes(ctx, tx, recipes.CreateRecipeNotesInput{
					RecipeID: recipeID,
					StepID:   step.StepID,
					Note:     step.Note,
				})
				if err != nil {
					return err
				}
			}
		}

		// Handle ingredient processing
		for _, ingredient := range input.Ingredients {
			err = recipes.CreateRecipeIngredients(ctx, tx, recipes.CreateRecipeIngredientInput{
				RecipeID:     recipeID,
				IngredientID: ingredient.IngredientID,
				Name:         ingredient.Name,
				Quantity:     ingredient.Quantity,
				Unit:         ingredient.Unit,
				Group:        ingredient.Group,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if txErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update recipe").SetInternal(err)
	}

	return c.JSON(http.StatusOK, recipe)
}

func (a *App) GetKitchenRecipe(c echo.Context) error {
	ctx := c.Request().Context()
	recipeID := c.Param("recipe_id")

	recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		if err == recipes.ErrRecipeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
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

	// Handle step hydration (images, notes)
	images, err := recipes.GetRecipeImagesByRecipeID(ctx, a.db, recipe.RecipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get images for recipe steps").SetInternal(err)
	}

	notes, err := recipes.GetRecipeNotesByRecipeID(ctx, a.db, recipe.RecipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get notes for recipes steps").SetInternal(err)
	}

	for i, step := range steps {
		for _, image := range images {
			if step.StepID == image.StepID {
				step.Images = append(step.Images, image.ImageURL)
			}
		}

		for _, note := range notes {
			if step.StepID == note.StepID {
				step.Note = note.Note
			}
		}

		steps[i] = step
	}

	recipe.Steps = steps

	ingredients, err := recipes.GetRecipeIngredientsByRecipeID(ctx, a.db, recipe.RecipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get ingredients for recipe").SetInternal(err)
	}

	recipe.Ingredients = ingredients

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

func (a *App) DeleteKitchenRecipe(c echo.Context) error {
	ctx := c.Request().Context()
	recipeID := c.Param("recipe_id")

	err := recipes.DeleteRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		if err == recipes.ErrRecipeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete recipe by ID").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

type SaveRecipeRequest struct {
	RecipeID string `json:"recipeId"`
}

func (a *App) SaveRecipe(c echo.Context) error {
	var input SaveRecipeRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	err = recipes.SaveRecipe(ctx, a.db, recipes.SaveRecipeInput{
		KitchenID: kitchenID,
		RecipeID:  input.RecipeID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not save recipe").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type RemoveSavedRecipeRequest struct {
	RecipeID string `json:"recipeId"`
}

func (a *App) RemoveSavedRecipe(c echo.Context) error {
	var input RemoveSavedRecipeRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	err = recipes.RemoveRecipe(ctx, a.db, recipes.RemoveRecipeInput{
		KitchenID: kitchenID,
		RecipeID:  input.RecipeID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not unsave recipe").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
