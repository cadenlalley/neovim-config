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
	err = json.Unmarshal(json.RawMessage(recipeText), &recipe)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}

func (a *App) ImportImage(c echo.Context) error {
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

	prefix := media.GetImportMediaPath(account.AccountID)
	key, err := a.handleFormFile(c, "file", prefix)
	if err != nil {
		if err == http.ErrMissingFile {
			return echo.NewHTTPError(http.StatusBadRequest, "no file provided")
		}
		err = errors.Wrapf(err, "could not upload file to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload file").SetInternal(err)
	}

	res, err := a.getRecipeFromImage(a.cdnHost + "/" + key)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not parse recipe from image").SetInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

func (a *App) getRecipeFromImage(url string) (recipes.Recipe, error) {
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
						Text: "Retrieve the complete recipe from the following image, including ingredient lists, quantities, and step-by-step instructions, preserving all original formatting and text.",
					},
					{
						Type: "image_url",
						ImageURL: &openai.ChatCompletionContentImageURL{
							URL: url,
						},
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
	err = json.Unmarshal(json.RawMessage(recipeText), &recipe)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}
