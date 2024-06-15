package api

import (
	"fmt"
	"net/http"

	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
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

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	prefix := fmt.Sprintf("uploads/%s/", account.AccountID)

	key, err := a.fileManager.UploadFromHeader(ctx, file, prefix)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload file").SetInternal(err)
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Key: key,
	})
}
