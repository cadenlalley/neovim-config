package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type CreateAccountRequest struct {
	UserID    string `json:"userId" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Kitchen   struct {
		Name   string `json:"name" validate:"required"`
		Bio    string `json:"bio"`
		Handle string `json:"handle" validate:"required"`
		Avatar string `json:"avatar"`
		Cover  string `json:"cover"`
		Public bool   `json:"public" validate:"required"`
	} `json:"kitchen" validate:"required"`
}

type CreateAccountResponse struct {
	Account accounts.Account `json:"account"`
	Kitchen kitchens.Kitchen `json:"kitchen"`
}

func (a *App) CreateAccount(c echo.Context) error {
	var input CreateAccountRequest
	if err := web.ValidateRequest(c, &input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userID := c.Get(auth.UserIDContextKey).(string)

	// Verify that the account being created actually belongs to the token making the request.
	if userID != input.UserID {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("userID does not match bearer"))
	}

	ctx := c.Request().Context()

	// Verify that the account does not already exist.
	account, err := accounts.GetAccountByUserID(ctx, a.db, input.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account").SetInternal(err)
	}

	if account.Exists() {
		return echo.NewHTTPError(http.StatusBadRequest, "account already exists")
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not start transaction").SetInternal(err)
	}
	defer tx.Rollback()

	account, err = accounts.CreateAccount(ctx, tx, accounts.CreateAccountInput{
		UserID:    input.UserID,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create account").SetInternal(err)
	}

	kitchen, err := kitchens.CreateKitchen(ctx, tx, kitchens.CreateKitchenInput{
		AccountID:   account.AccountID,
		KitchenName: input.Kitchen.Name,
		Bio:         input.Kitchen.Bio,
		Handle:      input.Kitchen.Handle,
		Avatar:      input.Kitchen.Avatar,
		Cover:       input.Kitchen.Cover,
		Public:      input.Kitchen.Public,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create kitchen").SetInternal(err)
	}

	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not commit account").SetInternal(err)
	}

	return c.JSON(http.StatusOK, CreateAccountResponse{
		Account: account,
		Kitchen: kitchen,
	})
}
