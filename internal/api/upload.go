package api

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
		if err == accounts.ErrAccountNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get account by user ID").SetInternal(err)
	}

	prefix := media.GetAccountMediaPath(account.AccountID)
	key, err := a.handleFormFile(c, "file", prefix)
	if err != nil {
		if err == http.ErrMissingFile {
			return echo.NewHTTPError(http.StatusBadRequest, "no file provided")
		}
		err = errors.Wrapf(err, "could not upload file to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload file").SetInternal(err)
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Key: key,
	})
}

func (a *App) CDN(c echo.Context) error {
	if a.env != ENV_DEV {
		return echo.NewHTTPError(http.StatusBadRequest, "not implemented")
	}

	ctx := c.Request().Context()
	parts := strings.Split(c.Request().URL.Path, "/cdn/")
	key := parts[len(parts)-1]

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no file requested")
	}

	file, contentType, err := a.fileManager.Get(ctx, key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Stream(http.StatusOK, contentType, file)
}

func (a *App) handleFormFile(c echo.Context, field, prefix string) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", err
	}

	ctx := c.Request().Context()
	keys, err := a.fileManager.UploadFromHeaders(ctx, []*multipart.FileHeader{file}, prefix)
	if err != nil {
		return "", err
	}

	return keys[0], nil
}

func (a *App) handleFormFiles(c echo.Context, fields []string, prefix string) ([]string, error) {
	files := make([]*multipart.FileHeader, 0)

	for _, field := range fields {
		file, err := c.FormFile(field)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	ctx := c.Request().Context()
	keys, err := a.fileManager.UploadFromHeaders(ctx, files, prefix)
	if err != nil {
		return nil, err
	}

	return keys, nil
}
