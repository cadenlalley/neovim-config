package middleware

import (
	"fmt"
	"net/http"

	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
)

type adminAuthorizer struct {
	UserIDs []string
}

func NewAdminAuthorizer(UserIDs []string) *adminAuthorizer {
	return &adminAuthorizer{
		UserIDs: UserIDs,
	}
}

// Validate that the user is an admin.
func (a *adminAuthorizer) ValidateAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get(auth.UserIDContextKey).(string)

		for _, id := range a.UserIDs {
			if id == userID {
				return next(c)
			}
		}

		err := fmt.Errorf("userID '%s' attempted to access admin resource without authorization", userID)
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}
}
