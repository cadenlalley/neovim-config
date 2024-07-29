package api

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
)

type GetIAMResponse struct {
	Account  accounts.Account   `json:"account"`
	Kitchens []kitchens.Kitchen `json:"kitchens"`
}

// GetIAM returns the profile for the user associated with the JWT.
func (a *App) GetIAM(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)

	// Lookup the user record for the provided JWT.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account by user ID").SetInternal(err)
	}

	// If the account does not exist, we can populate it partially for the frontend to finish out during setup.
	if account.AccountID == "" {
		claims := c.Get(auth.ClaimsContextKey).(*validator.ValidatedClaims).CustomClaims.(*auth.CustomClaims)

		return c.JSON(http.StatusOK, GetIAMResponse{
			Account: accounts.Account{
				UserID:    userID,
				Email:     claims.Email,
				FirstName: claims.FirstName,
				LastName:  claims.LastName,
			},
		})
	}

	// Find all kitchens for the provided account.
	accountKitchens, err := kitchens.ListKitchensByAccountID(ctx, a.db, account.AccountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get kitchen(s)by account ID")
	}

	// If the account exists, look up the associated profiles.
	return c.JSON(http.StatusOK, GetIAMResponse{
		Account:  account,
		Kitchens: accountKitchens,
	})
}
