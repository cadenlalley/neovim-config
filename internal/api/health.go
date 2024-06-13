package api

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

// Get tests the health of the application.
func (a *App) GetHealth(c echo.Context) error {
	ctx := c.Request().Context()

	if err := a.db.Ping(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "database unhealthy").SetInternal(err)
	}

	_, err := a.s3.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "object store unhealthy").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
