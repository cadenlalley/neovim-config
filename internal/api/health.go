package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get tests the health of the application.
func (a *App) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()

	if err := a.db.Ping(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database unhealthy").SetInternal(err)
	}

	if err := a.fileManager.Ping(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "file store unhealthy").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
