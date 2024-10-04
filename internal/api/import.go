package api

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/extractor"
	"github.com/kitchens-io/kitchens-api/internal/openai"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
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
		ResponseFormat: json.RawMessage(`{
        "type": "json_schema",
        "json_schema": {
            "name": "url_recipe_response",
            "schema": {
                "type": "object",
                "properties": {
                    "name": {"type": "string"},
                    "summary": {"type": "string"},
                    "prepTime": {"type": "integer", "description": "the preparation time in minutes"},
                    "cookTime": {"type": "integer", "description": "the cook time in minutes"},
                    "servings": {"type": "integer", "description": "the number of servings"},
                    "ingredients": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "ingredientId": {"type": "integer"},
                                "name": {"type": "string"},
                                "quantity": {"type": "number"},
                                "unit": {"type": ["string", "null"], "description": "the unit of measurement"},
                                "group": {"type": ["string", "null"]}
                            },
                            "required": ["ingredientId","name","quantity","group","unit"],
                            "additionalProperties": false
                        }
                    },
                    "steps": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "stepId": {"type": "integer"},
                                "instruction": {"type": "string"},
                                "note": {"type": ["string", "null"]},
                                "group": {"type": ["string", "null"]}
                            },
                            "required": ["stepId","instruction","note","group"],
                            "additionalProperties": false
                        }
                    }
                },
                "required": ["name","summary","prepTime","cookTime","servings","ingredients","steps"],
                "additionalProperties": false
            },
            "strict": true
        }
    }`),
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
