package ai

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var RecipeResponseJSONSchema = GenerateSchema[RecipeResponseSchema]()

type RecipeResponseSchema struct {
	Name        string                            `json:"name"`
	Summary     string                            `json:"summary"`
	PrepTime    int                               `json:"prepTime" jsonschema_description:"the prep time in minutes"`
	CookTime    int                               `json:"cookTime" jsonschema_description:"the cook time in minutes"`
	Servings    int                               `json:"servings" jsonschema_description:"the number of servings"`
	Ingredients []RecipeResponseIngredientsSchema `json:"ingredients"`
	Steps       []RecipeResponseStepsSchema       `json:"steps"`
}

type RecipeResponseIngredientsSchema struct {
	IngredientID int     `json:"ingredientId"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity" jsonschema_description:"the amount of the ingredient, if an amount doesn't make sense for the ingredient set it to 0"`
	Unit         string  `json:"unit" jsonschema:"enum=bag,enum=bottle,enum=box,enum=can,enum=clove,enum=cup,enum=dash,enum=drop,enum=gallon,enum=gram,enum=jar,enum=kilogram,enum=liter,enum=milliliter,enum=ounce,enum=packet,enum=piece,enum=pint,enum=pinch,enum=pound,enum=quart,enum=slice,enum=stick,enum=tbsp,enum=tsp,enum=n/a" jsonschema_description:"optional unit of measurement, if a unit doesn't make sense for the ingredient set it to n/a"`
	Group        string  `json:"group" jsonschema_description:"the group within the recipe the ingredient belongs to"`
}

type RecipeResponseStepsSchema struct {
	StepID      int    `json:"stepId"`
	Instruction string `json:"instruction"`
	Note        string `json:"note"`
	Group       string `json:"group" jsonschema_description:"the group within the recipe the step belongs to"`
}

func (a *AIClient) ExtractRecipeFromText(ctx context.Context, text string) (RecipeResponseSchema, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "recipe",
		Description: openai.String("A JSON object representing a recipe"),
		Schema:      RecipeResponseJSONSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModelGPT4oMini,
		MaxTokens: openai.Int(1600),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Retrieve the complete recipe from the following text, including ingredient lists, quantities, and step-by-step instructions, preserving all original formatting and text: " + text),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return RecipeResponseSchema{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	log.Info().
		Str("producer", "openai_text_import").
		Interface("tokenUsage", TokenUsage{
			PromptTokens:     chat.Usage.PromptTokens,
			CompletionTokens: chat.Usage.CompletionTokens,
			TotalTokens:      chat.Usage.TotalTokens,
		}).
		Msg("openai metadata")

	if chat.Choices == nil || len(chat.Choices) == 0 {
		return RecipeResponseSchema{}, errors.New("no choices returned from OpenAI API")
	}

	var results RecipeResponseSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return RecipeResponseSchema{}, errors.Wrap(err, "failed to unmarshal recipe response")
	}

	return results, nil
}

func (a *AIClient) ExtractRecipeFromImageURLs(ctx context.Context, urls []string) (RecipeResponseSchema, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "image_recipe_response",
		Description: openai.String("A JSON object representing a recipe extracted from a URL."),
		Schema:      RecipeResponseJSONSchema,
		Strict:      openai.Bool(true),
	}

	// Append each URL that has been provided to completion.
	images := []openai.ChatCompletionContentPartUnionParam{}
	for _, url := range urls {
		images = append(images, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL: url,
		}))
	}

	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModelGPT4oMini,
		MaxTokens: openai.Int(1600),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Retrieve the complete recipe from the following images, including ingredient lists, quantities, and step-by-step instructions, preserving all original formatting and text."),
			openai.UserMessage(images),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return RecipeResponseSchema{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	log.Info().
		Str("producer", "openai_image_import").
		Interface("tokenUsage", TokenUsage{
			PromptTokens:     chat.Usage.PromptTokens,
			CompletionTokens: chat.Usage.CompletionTokens,
			TotalTokens:      chat.Usage.TotalTokens,
		}).
		Msg("openai metadata")

	if len(chat.Choices) == 0 {
		return RecipeResponseSchema{}, errors.New("no choices returned from OpenAI API")
	}

	var results RecipeResponseSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return RecipeResponseSchema{}, errors.Wrap(err, "failed to unmarshal recipe response")
	}

	return results, nil
}
