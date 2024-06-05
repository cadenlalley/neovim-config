package api

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/models"
	"github.com/labstack/echo/v4"
)

type GetIAMResponse struct {
	Account  models.Account   `json:"account"`
	Profiles []models.Profile `json:"profiles"`
}

// GetIAM returns the profile for the user associated with the JWT.
func (a *App) GetIAM(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)

	// Lookup the user record for the provided JWT.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// If the account does not exist, we can populate it partially for the frontend to finish out during setup.
	if account.AccountID == "" {
		claims := c.Get(auth.ClaimsContextKey).(*validator.ValidatedClaims).CustomClaims.(*auth.CustomClaims)

		return c.JSON(http.StatusOK, GetIAMResponse{
			Account: models.Account{
				UserID:    userID,
				Email:     claims.Email,
				FirstName: claims.FirstName,
				LastName:  claims.LastName,
			},
		})
	}

	// If the account exists, look up the associated profiles.
	return c.JSON(http.StatusOK, GetIAMResponse{
		Account: account,
	})
}
