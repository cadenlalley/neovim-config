package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/models"
	"github.com/labstack/echo/v4"
)

type CreateAccountRequest struct {
	UserID    string `json:"userId" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type CreateAccountResponse struct {
	Account models.Account `json:"account"`
}

func (a *App) CreateAccount(c echo.Context) error {
	var input CreateAccountRequest
	if err := web.ValidateRequest(c, &input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	// Verify that the account does not already exist.
	account, err := accounts.GetAccountByUserID(ctx, a.db, input.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if account.Exists() {
		return c.JSON(http.StatusOK, CreateAccountResponse{
			Account: account,
		})
	}

	account, err = accounts.CreateAccount(ctx, a.db, accounts.CreateAccountInput{
		UserID:    input.UserID,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, CreateAccountResponse{
		Account: account,
	})
}
