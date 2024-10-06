package api

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/override"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type CreateAccountRequest struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Kitchen   struct {
		Name    string `form:"kitchenName" validate:"required"`
		Bio     string `form:"kitchenBio"`
		Handle  string `form:"kitchenHandle" validate:"required"`
		Private bool   `form:"kitchenPrivate"`

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
	if err != nil && err != accounts.ErrAccountNotFound {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account").SetInternal(err)
	}

	if account.Exists() {
		return echo.NewHTTPError(http.StatusBadRequest, accounts.ErrDuplicateAccount)
	}

	// Create an account ID.
	accountID := accounts.CreateAccountID()
	kitchenID := kitchens.CreateKitchenID()

	// Handle the file uploads if they have been set.
	var kitchen kitchens.Kitchen

	prefix := media.GetKitchenMediaPath(kitchenID)
	avatarKey, err := a.handleFormFile(c, "kitchenAvatarFile", prefix)
	if err != nil && err != http.ErrMissingFile {
		err = errors.Wrapf(err, "could not upload kitchenAvatarFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload avatar photo").SetInternal(err)
	} else {
		kitchen.Avatar = null.NewString(avatarKey, true)
	}

	coverKey, err := a.handleFormFile(c, "kitchenCoverFile", prefix)
	if err != nil && err != http.ErrMissingFile {
		err = errors.Wrapf(err, "could not upload kitchenCoverFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload cover photo").SetInternal(err)
	} else {
		kitchen.Cover = null.NewString(coverKey, true)
	}

	err = mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		account, err = accounts.CreateAccount(ctx, tx, accounts.CreateAccountInput{
			AccountID: accountID,
			UserID:    userID,
			Email:     claims.Email,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		})
		if err != nil {
			return err
		}

		kitchen, err = kitchens.CreateKitchen(ctx, tx, kitchens.CreateKitchenInput{
			KitchenID:   kitchenID,
			AccountID:   account.AccountID,
			KitchenName: input.Kitchen.Name,
			Bio:         input.Kitchen.Bio,
			Handle:      input.Kitchen.Handle,
			Avatar:      avatarKey,
			Cover:       coverKey,
			Private:     input.Kitchen.Private,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == kitchens.ErrDuplicateHandle {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create account").SetInternal(err)
	}

	return c.JSON(http.StatusOK, CreateAccountResponse{
		Account: account,
		Kitchen: kitchen,
	})
}

type UpdateAccountRequest struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
}

func (a *App) UpdateAccount(c echo.Context) error {
	var input UpdateAccountRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)

	// Verify that the account does not already exist.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		if err == accounts.ErrAccountNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account").SetInternal(err)
	}

	account, err = accounts.UpdateAccount(ctx, a.db, accounts.UpdateAccountInput{
		AccountID: account.AccountID,
		FirstName: override.String(input.FirstName, account.FirstName),
		LastName:  override.String(input.LastName, account.LastName),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update account").SetInternal(err)
	}

	return c.JSON(http.StatusOK, account)
}
