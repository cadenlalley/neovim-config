package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
)

func (a *App) GetKitchenRecipe(c echo.Context) error {
	ctx := c.Request().Context()
	recipeID := c.Param("recipe_id")

	recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe by ID").SetInternal(err)
	}

	// TODO: Add fetching recipe dependencies.

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
