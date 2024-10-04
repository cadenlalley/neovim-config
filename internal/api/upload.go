package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UploadResponse struct {
	Key string `json:"key"`
}

func (a *App) Upload(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get(auth.UserIDContextKey).(string)

	// Lookup the user record for the provided JWT.
	account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account by user ID").SetInternal(err)
	}

	prefix := media.GetAccountMediaPath(account.AccountID)
	key, err := a.HandleFormFile(c, "file", prefix)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no file provided") {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		msg := "could not upload file"
		log.Err(err).Str("prefix", prefix).Msg(msg)
		return echo.NewHTTPError(http.StatusInternalServerError, msg).SetInternal(err)
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Key: key,
	})
}

func (a *App) HandleFormFile(c echo.Context, field, prefix string) (string, error) {
	file, err := c.FormFile(field)
	if file != nil && err != nil {
		return "", err
	}

	if file == nil {
		return "", fmt.Errorf("no file provided in field '%s'", field)
	}

	ctx := c.Request().Context()
	key, err := a.fileManager.UploadFromHeader(ctx, file, prefix)
	if err != nil {
		return "", err
	}

	return key, nil
}
