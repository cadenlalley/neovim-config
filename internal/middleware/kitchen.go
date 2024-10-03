package middleware

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
)

type KitchenAuthorizer struct {
	db *sqlx.DB
}

func NewKitchenAuthorizer(db *sqlx.DB) *KitchenAuthorizer {
	return &KitchenAuthorizer{
		db: db,
	}
}

// Validate that the user has permission to this kitchen.
func (a *KitchenAuthorizer) ValidateWriter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := c.Get(auth.UserIDContextKey).(string)
		kitchenID := c.Param("kitchen_id")

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
			err := fmt.Errorf("account '%s' attempted to create in kitchen '%s' without authorization", account.AccountID, kitchen.KitchenID)
			return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
		}

		return next(c)
	}
}
