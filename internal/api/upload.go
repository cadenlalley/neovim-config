package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/ptr"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
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

	parts := strings.Split(file.Filename, ".")
	fileType := parts[len(parts)-1]

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not open file").SetInternal(err)
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read file").SetInternal(err)
	}

	key := fmt.Sprintf("uploads/%s/%s.%s", account.AccountID, ksuid.New().String(), fileType)
	_, err = a.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: ptr.String("kitchens-app-local-us-east-1"),
		Key:    ptr.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload file").SetInternal(err)
	}

	return c.JSON(http.StatusOK, UploadResponse{
		Key: key,
	})
}
