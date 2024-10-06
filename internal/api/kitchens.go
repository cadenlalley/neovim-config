package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/override"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
	kitchenID := c.Param("kitchen_id")

	var input UpdateKitchenRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	kitchen, err := kitchens.GetKitchenByID(ctx, a.db, kitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get kitchen by ID").SetInternal(err)
	}

	// Since bio is allowed to be null
	bio := kitchen.Bio.String
	if input.Bio != nil {
		bio = *input.Bio
	}

	// Handle the file uploads if they have been set.
	prefix := media.GetKitchenMediaPath(kitchenID)
	avatarKey, err := a.handleFormFile(c, "avatarFile", prefix)
	if err != nil && err != http.ErrMissingFile {
		err = errors.Wrapf(err, "could not upload avatarFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload avatar photo").SetInternal(err)
	}

	// Since avatarFile is allowed to be null.
	avatarFile := kitchen.Avatar.String
	if avatarKey != "" {
		avatarFile = avatarKey
	}

	coverKey, err := a.handleFormFile(c, "coverFile", prefix)
	if err != nil && err != http.ErrMissingFile {
		err = errors.Wrapf(err, "could not upload coverFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload cover photo").SetInternal(err)
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
