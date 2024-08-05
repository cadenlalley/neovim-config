package api

import (
	"fmt"
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/override"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	Name    string  `form:"name"`
	Bio     *string `form:"bio"`
	Handle  string  `form:"handle"`
	Private *bool   `form:"private"`

	// The following are manually checked in the CreateAccount handler.
	// They cannot be bound automatically, and are optional.
	//
	// kitchenAvatarFile, kitchenCoverFile
	AvatarFile  string
	KitchenFile string
}

func (a *App) UpdateKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)
	kitchenID := c.Param("kitchen_id")

	var input UpdateKitchenRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
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

	// Since bio is allowed to be null
	bio := kitchen.Bio.String
	if input.Bio != nil {
		bio = *input.Bio
	}

	// Handle the file uploads if they have been set.
	prefix := media.GetKitchenMediaPath(kitchenID)
	avatarKey, err := a.HandleFormFile(c, "avatarFile", prefix)
	if err != nil {
		msg := "could not upload avatar photo"
		log.Err(err).Str("prefix", prefix).Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg).SetInternal(err)
	}

	// Since avatarFile is allowed to be null.
	avatarFile := kitchen.Avatar.String
	if avatarKey != "" {
		avatarFile = avatarKey
	}

	coverKey, err := a.HandleFormFile(c, "coverFile", prefix)
	if err != nil {
		msg := "could not upload cover photo"
		log.Err(err).Str("prefix", prefix).Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg).SetInternal(err)
	}

	// Since coverFile is allowed to be null.
	coverFile := kitchen.Cover.String
	if coverKey != "" {
		coverFile = coverKey
	}

	kitchen, err = kitchens.UpdateKitchen(ctx, a.db, kitchens.UpdateKitchenInput{
		KitchenID: kitchenID,
		Name:      override.String(input.Name, kitchen.Name),
		Bio:       override.NullString(bio, kitchen.Bio),
		Handle:    override.String(input.Handle, kitchen.Handle),
		Avatar:    override.NullString(avatarFile, kitchen.Avatar),
		Cover:     override.NullString(coverFile, kitchen.Cover),
		Private:   override.Bool(input.Private, kitchen.Private),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update kitchen").SetInternal(err)
	}

	return c.JSON(http.StatusOK, kitchen)
}
