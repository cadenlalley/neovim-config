package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/labstack/echo/v4"
	"github.com/openai/openai-go"
	oaiOpt "github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/pkg/errors"
)

// TODO: No idea why this needs to be added to every request,
// but if it's set on the client then there are issues with values being passed around.
func getOpenAIOpts() []oaiOpt.RequestOption {
	var openAIClientOps = []oaiOpt.RequestOption{
		oaiOpt.WithAPIKey(os.Getenv("OPENAI_TOKEN")),
		oaiOpt.WithBaseURL("https://api.openai.com/v1"),
	}
	return openAIClientOps
}

type ExternalSearchResults struct {
	Recipes []ExternalSearchResult `json:"recipes"`
}

type ExternalSearchResult struct {
	Name         string  `json:"name" jsonschema_description:"Name of the recipe"`
	Source       string  `json:"source" jsonschema_description:"URL source of the recipe"`
	ReviewRating float64 `json:"reviewRating" jsonschema_description:"Rating of the recipe, 0 if not specified"`
	ReviewCount  int     `json:"reviewCount" jsonschema_description:"Number of reviews for the recipe, 0 if not specified"`
}

var externalSearchResultsSchema = GenerateSchema[ExternalSearchResults]()

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func (a *App) ExternalSearch(c echo.Context) error {
	ctx := c.Request().Context()
	query := c.QueryParam("q")

	// Ask OpenAI to perform an external search based on the provided query.
	results, err := a.getExternalSearchResults(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not perform external search").SetInternal(err)
	}

	return c.JSON(http.StatusOK, results)
}

func (a *App) getExternalSearchResults(ctx context.Context, query string) (ExternalSearchResults, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "recipe_search_results",
		Description: openai.String("A list of recipes that match the search query"),
		Schema:      externalSearchResultsSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := a.aiClientV2.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModelGPT4oMiniSearchPreview,
		MaxTokens: param.Opt[int64]{Value: 1600},
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
	}, getOpenAIOpts()...)
	if err != nil {
		return ExternalSearchResults{}, errors.Wrap(err, "failed to get external search results")
	}

	var results ExternalSearchResults
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return ExternalSearchResults{}, errors.Wrap(err, "failed to unmarshal external search results")
	}

	return results, nil
}
