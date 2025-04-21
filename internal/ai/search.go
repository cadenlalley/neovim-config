package ai

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var ExternalSearchResultsJSONSchema = GenerateSchema[ExternalSearchResults]()

type ExternalSearchResults struct {
	Recipes []ExternalSearchResult `json:"recipes"`
}

type ExternalSearchResult struct {
	Name         string  `json:"name" jsonschema_description:"Name of the recipe"`
	Source       string  `json:"source" jsonschema_description:"URL source of the recipe"`
	ReviewRating float64 `json:"reviewRating" jsonschema_description:"Rating of the recipe, 0 if not specified"`
	ReviewCount  int     `json:"reviewCount" jsonschema_description:"Number of reviews for the recipe, 0 if not specified"`
}

func (a *AIClient) ExternalSearch(ctx context.Context, query string) (ExternalSearchResults, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "recipe_search_results",
		Description: openai.String("A list of recipes that match the search query"),
		Schema:      ExternalSearchResultsJSONSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModelGPT4oMiniSearchPreview,
		MaxTokens: openai.Int(1600),
		WebSearchOptions: openai.ChatCompletionNewParamsWebSearchOptions{
			SearchContextSize: "high",
		},
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Return 10 " + query + " recipes. Prefer recipes from reputable websites and that have high ratings and high review counts. Prefer a variety of domains. Prefer webpages that are a single recipe. Each recipe should include a direct URL to the original source, the average rating, and total reviews."),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})
	if err != nil {
		return ExternalSearchResults{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	log.Info().
		Str("producer", "openai_external_search").
		Interface("tokenUsage", TokenUsage{
			PromptTokens:     chat.Usage.PromptTokens,
			CompletionTokens: chat.Usage.CompletionTokens,
			TotalTokens:      chat.Usage.TotalTokens,
		}).
		Msg("openai metadata")

	if chat.Choices == nil || len(chat.Choices) == 0 {
		return ExternalSearchResults{}, errors.New("no choices returned from OpenAI API")
	}

	var results ExternalSearchResults
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return ExternalSearchResults{}, errors.Wrap(err, "failed to unmarshal external search results")
	}

	return results, nil
}
