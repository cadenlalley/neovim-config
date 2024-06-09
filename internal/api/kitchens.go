package api

import (
	"fmt"
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
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

type UpdateKitchenRequest struct {
	Name   string `json:"name" validate:"required"`
	Bio    string `json:"bio" validate:"required"`
	Handle string `json:"handle" validate:"required"`
	Avatar string `json:"avatar" validate:"required"`
	Cover  string `json:"cover" validate:"required"`
	Public *bool  `json:"public" validate:"required"`
}

func (a *App) UpdateKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)
	kitchenID := c.Param("kitchen_id")

	var input UpdateKitchenRequest
	if err := web.ValidateRequest(c, &input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

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
		err := fmt.Errorf("account '%s' attempted to modify kitchen '%s' without authorization", account.AccountID, kitchen.KitchenID)
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}

	// Convert bool to value
	public := true
	if input.Public != nil && !*input.Public {
		public = false
	}

	kitchen, err = kitchens.UpdateKitchen(ctx, a.db, kitchens.UpdateKitchenInput{
		KitchenID:   kitchenID,
		KitchenName: input.Name,
		Bio:         input.Bio,
		Handle:      input.Handle,
		Avatar:      input.Avatar,
		Cover:       input.Cover,
		Public:      public,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update kitchen").SetInternal(err)
	}

	return c.JSON(http.StatusOK, kitchen)
}
