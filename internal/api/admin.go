package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *App) AdminListAccounts(c echo.Context) error {
	ctx := c.Request().Context()

	accounts, err := accounts.ListAccounts(ctx, a.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not list accounts").SetInternal(err)
	}

	return c.JSON(http.StatusOK, accounts)
}

type CreateRecipeMetadataRequest struct {
	RecipeID string `json:"recipeId"`
}

func (a *App) AdminCreateRecipeMetadata(c echo.Context) error {
	var input CreateRecipeMetadataRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	result, err := a.extractRecipeMeta(ctx, input.RecipeID)
	if err != nil {
		if errors.Is(err, recipes.ErrRecipeNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not extract recipe metadata").SetInternal(err)
	}

	return c.JSON(http.StatusOK, result)
}
