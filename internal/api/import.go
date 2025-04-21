package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kitchens-io/kitchens-api/internal/extractor"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type ImportURLRequest struct {
	Source   string `json:"source"`
	Features struct {
		Groups bool `json:"groups"`
	} `json:"features"`
	Debug struct {
		Prompt string `json:"prompt"`
	} `json:"debug"`
}

func (a *App) ImportURL(c echo.Context) error {
	var input ImportURLRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()

	// Load the requested source
	str, err := extractor.GetTextFromURL(input.Source)
	if err != nil {
		if err == extractor.ErrRequestBlocked {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	res, err := a.aiClient.ExtractRecipeFromText(ctx, str)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	data, err := json.Marshal(res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	var r recipes.Recipe
	err = r.Import(data, input.Features.Groups)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}
	r.Source = null.NewString(input.Source, true)

	return c.JSON(http.StatusOK, r)
}

type ImportImageRequest struct {
	Count    int `form:"count" validate:"required"`
	Features struct {
		Groups bool `form:"featureGroups"`
	}
	URLs string `form:"urls"`

	// The following are manually checked in the handler based on the provided count.
	// file_1, file_2, file_n...
}

func (a *App) ImportImage(c echo.Context) error {
	var input ImportImageRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

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

	// Based on the count provided, generate the number of files there should be.
	fields := make([]string, 0)
	for i := 1; i < input.Count+1; i++ {
		fields = append(fields, fmt.Sprintf("file_%d", i))
	}

	prefix := media.GetImportMediaPath(account.AccountID)

	keys, err := a.handleFormFiles(c, fields, prefix)
	if err != nil {
		if err == http.ErrMissingFile {
			return echo.NewHTTPError(http.StatusBadRequest, "no file provided")
		}
		err = errors.Wrapf(err, "could not upload file to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload file").SetInternal(err)
	}

	urls := make([]string, 0)
	for _, key := range keys {
		urls = append(urls, a.cdnHost+"/"+key)
	}

	// Image Uploads don't work in development, however we can return an empty recipe for debugging.
	// If URLs are provided in the form submission, use those instead.
	if a.env == ENV_DEV {
		if len(input.URLs) > 0 {
			parts := strings.Split(input.URLs, ",")
			urls = make([]string, 0)
			urls = append(urls, parts...)
		} else {
			var sample recipes.Recipe
			err := json.Unmarshal(recipes.Sample, &sample)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not parse sample recipe for development").SetInternal(err)
			}
			return c.JSON(http.StatusOK, sample)
		}
	}

	res, err := a.aiClient.ExtractRecipeFromImageURLs(ctx, urls)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	data, err := json.Marshal(res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	var r recipes.Recipe
	err = r.Import(data, input.Features.Groups)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
