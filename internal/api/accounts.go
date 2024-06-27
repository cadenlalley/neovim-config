package api

import (
	"fmt"
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

type CreateAccountRequest struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Kitchen   struct {
		Name   string `form:"kitchenName" validate:"required"`
		Bio    string `form:"kitchenBio"`
		Handle string `form:"kitchenHandle" validate:"required"`
		Public bool   `form:"kitchenPublic" validate:"required"`

		// The following are manually checked in the CreateAccount handler.
		// They cannot be bound automatically, and are optional.
		//
		// kitchenAvatarFile, kitchenCoverFile
	}
}

type CreateAccountResponse struct {
	Account accounts.Account `json:"account"`
	Kitchen kitchens.Kitchen `json:"kitchen"`
}

func (a *App) CreateAccount(c echo.Context) error {
	var input CreateAccountRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userID := c.Get(auth.UserIDContextKey).(string)
	claims := c.Get(auth.ClaimsContextKey).(*validator.ValidatedClaims).CustomClaims.(*auth.CustomClaims)

	ctx := c.Request().Context()

	// Verify that the account does not already exist.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account").SetInternal(err)
	}

	if account.Exists() {
		return echo.NewHTTPError(http.StatusBadRequest, "account already exists")
	}

	// Handle the file uploads if they have been set.
	var kitchen kitchens.Kitchen

	prefix := fmt.Sprintf("uploads/%s/", userID)

	avatarKey, err := a.HandleFormFile(c, "kitchenAvatarFile", prefix)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload avatar photo").SetInternal(err)
	} else {
		kitchen.Avatar = null.NewString(avatarKey, true)
	}

	coverKey, err := a.HandleFormFile(c, "kitchenCoverFile", prefix)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload cover photo").SetInternal(err)
	} else {
		kitchen.Cover = null.NewString(coverKey, true)
	}

	err = mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		account, err = accounts.CreateAccount(ctx, tx, accounts.CreateAccountInput{
			UserID:    userID,
			Email:     claims.Email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		})
		if err != nil {
			return err
		}

		kitchen, err = kitchens.CreateKitchen(ctx, tx, kitchens.CreateKitchenInput{
			AccountID:   account.AccountID,
			KitchenName: input.Kitchen.Name,
			Bio:         input.Kitchen.Bio,
			Handle:      input.Kitchen.Handle,
			Avatar:      avatarKey,
			Cover:       coverKey,
			Public:      input.Kitchen.Public,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create account").SetInternal(err)
	}

	return c.JSON(http.StatusOK, CreateAccountResponse{
		Account: account,
		Kitchen: kitchen,
	})
}
