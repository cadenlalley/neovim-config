package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
)

func (a *App) GetKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	kitchen, err := kitchens.GetKitchenByID(ctx, a.db, kitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get kitchen by ID").SetInternal(err)
	}

	return c.JSON(http.StatusOK, kitchen)
}
