package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) ExternalSearch(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")

	// Ask OpenAI to perform an external search based on the provided query.
	results, err := a.aiClient.ExternalSearch(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not perform external search").SetInternal(err)
	}

	return c.JSON(http.StatusOK, results)
}
