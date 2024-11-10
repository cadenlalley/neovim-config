package api

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/extractor"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/openai"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type ImportURLRequest struct {
	Source string `json:"source"`
}

func (a *App) ImportURL(c echo.Context) error {
	var input ImportURLRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Load the requested source
	str, err := extractor.GetTextFromURL(input.Source)
	if err != nil {
		if err == extractor.ErrRequestBlocked {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	// Ask OpenAI to convert the string to a structured recipe.
	res, err := a.getRecipeFromText(str)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from URL").SetInternal(err)
	}

	res.Source = null.NewString(input.Source, true)

	return c.JSON(http.StatusOK, res)
}

func (a *App) getRecipeFromText(text string) (recipes.Recipe, error) {
	res, err := a.aiClient.PostChatCompletion(openai.ChatCompletionRequest{
		Model:     "gpt-4o-mini",
		MaxTokens: 1600,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: []openai.ChatCompletionContent{
					{
						Type: "text",
						Text: "You are a helpful assistant.",
					},
				},
			},
			{
				Role: "user",
				Content: []openai.ChatCompletionContent{
					{
						Type: "text",
						Text: "Retrieve the complete recipe from the following text, including ingredient lists, quantities, and step-by-step instructions, preserving all original formatting and text: " + text,
					},
				},
			},
		},
		ResponseFormat: recipes.JsonSchema,
	})
	if err != nil {
		return recipes.Recipe{}, err
	}

	// Handle the escaped HTML that comes back from Open AI.
	if len(res.Choices) == 0 {
		return recipes.Recipe{}, fmt.Errorf("unexpected response payload from Open AI: %v", res)
	}

	recipeText := html.UnescapeString(res.Choices[0].Message.Content)

	// Marshal the response into a recipe struct.
	var recipe recipes.Recipe
	err = recipe.Import(json.RawMessage(recipeText), true)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}

type ImportImageRequest struct {
	Count int `form:"count" validate:"required"`

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
	if a.env == ENV_DEV {
		var sample recipes.Recipe
		err := json.Unmarshal(recipes.Sample, &sample)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not parse sample recipe for development").SetInternal(err)
		}
		return c.JSON(http.StatusOK, sample)
	}

	res, err := a.getRecipeFromImages(urls)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from image").SetInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

func (a *App) getRecipeFromImages(urls []string) (recipes.Recipe, error) {

	chatCompletion := []openai.ChatCompletionContent{
		{
			Type: "text",
			Text: "Retrieve the complete recipe from the following images, including ingredient lists, quantities, and step-by-step instructions, preserving all original formatting and text.",
		},
	}

	// Append each URL that has been provided to completion.
	for _, url := range urls {
		chatCompletion = append(chatCompletion, openai.ChatCompletionContent{
			Type: "image_url",
			ImageURL: &openai.ChatCompletionContentImageURL{
				URL: url,
			},
		})
	}

	res, err := a.aiClient.PostChatCompletion(openai.ChatCompletionRequest{
		Model:     "gpt-4o-mini",
		MaxTokens: 1600,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: []openai.ChatCompletionContent{
					{
						Type: "text",
						Text: "You are a helpful assistant.",
					},
				},
			},
			{
				Role:    "user",
				Content: chatCompletion,
			},
		},
		ResponseFormat: recipes.JsonSchema,
	})
	if err != nil {
		return recipes.Recipe{}, err
	}

	// Handle the escaped HTML that comes back from Open AI.
	if len(res.Choices) == 0 {
		return recipes.Recipe{}, fmt.Errorf("unexpected response payload from Open AI: %v", res)
	}

	recipeText := html.UnescapeString(res.Choices[0].Message.Content)

	// Marshal the response into a recipe struct.
	var recipe recipes.Recipe
	err = recipe.Import(json.RawMessage(recipeText), false)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}
