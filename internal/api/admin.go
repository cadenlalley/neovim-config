package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/labstack/echo/v4"
)

func (a *App) AdminListAccounts(c echo.Context) error {
	ctx := c.Request().Context()

	accounts, err := accounts.ListAccounts(ctx, a.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not list accounts").SetInternal(err)
	}

	return c.JSON(http.StatusOK, accounts)
}
