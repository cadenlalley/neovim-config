package ai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/metrics"
	"github.com/openai/openai-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var RecipeMetaResponseJSONSchema = GenerateSchema[RecipeMetaResponseSchema]()

type RecipeMetaResponseSchema struct {
	Difficulty int    `json:"difficulty" jsonschema:"enum=1,enum=2,enum=3,enum=4,enum=5"`
	Course     string `json:"course" jsonschema:"enum=breakfast,enum=brunch,enum=lunch,enum=dinner,enum=dessert,enum=supper"`
	Class      string `json:"class" jsonschema:"enum=main,enum=side,enum=snack,enum=beverage,enum=dessert,enum=dip,enum=soup,enum=appetizer"`
	Cuisine    string `json:"cuisine" jsonschema_description:"the cuisine of the recipe"`
	Tags       []struct {
		Type  string `json:"type" jsonschema_description:"enum=diet,enum=keyword,enum=ingredient"`
		Value string `json:"value" jsonschema_description:"the value of the tag"`
	} `json:"tags" jsonschema_description:"an array of tags"`
}

func (a *AIClient) ExtractRecipeMetaFromText(ctx context.Context, text string) (RecipeMetaResponseSchema, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "recipe_meta",
		Description: openai.String("A JSON object representing tagging metadata for a recipe"),
		Schema:      RecipeMetaResponseJSONSchema,
		Strict:      openai.Bool(true),
	}

	start := time.Now()
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModelGPT4oMini,
		MaxTokens: openai.Int(1600),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(`
Your task is to parse the provided JSON recipe into search‑friendly tags.

**TAGGING RULES**
1. **difficulty**: one tag only, value "1"…"5" (stringified int; 1 = easiest).
2. **cuisine**: one tag only, the cuisine that best represents the recipe (e.g. italian, mexican, korean, mediterranean, tex-mex).
3. **course**: exactly one of: breakfast, brunch, lunch, dinner, dessert, supper.
4. **class**: exactly one of: main, side, snack, beverage, dessert, dip, soup, appetizer.
5. **diet**: one tag per diet/allergen (e.g. vegan, vegetarian, keto, paleo, high-protein, low-carb).
6. **ingredient**: 2–3 defining ingredients (singular nouns, e.g. chicken-breast, ground-beef, salmon-filet, black-bean).
7. **keyword**: 5–12 extra tokens (2-3 words each) that aid search.
   - Must **not duplicate** any value already tagged under other types.
   - Buckets to consider:
		- equipment/technique (e.g. instant-pot, air-fryer, one-pan, no-bake)
		- time/effort (under‑15‑min, 30-minute-meal, overnight)
		- occasion (meal‑prep)
		- audience (kid‑friendly, budget, beginner)

**FORMAT & STYLE**
* All strings must be lowercase; multi‑word values must use kebab‑case (sheet‑pan).
* Each tag object must include both "type" and "value".
* If a tag is not applicable, do not include it.

Follow these rules strictly to produce clean, indexable metadata for every recipe.

---

**Recipe JSON:**

` + text),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return RecipeMetaResponseSchema{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	log.Info().
		Str("producer", "openai_text_import").
		Interface("tokenUsage", TokenUsage{
			PromptTokens:     chat.Usage.PromptTokens,
			CompletionTokens: chat.Usage.CompletionTokens,
			TotalTokens:      chat.Usage.TotalTokens,
		}).
		Int64("latency", metrics.Elapsed(start)).
		Msg("openai metadata")

	if chat.Choices == nil || len(chat.Choices) == 0 {
		return RecipeMetaResponseSchema{}, errors.New("no choices returned from OpenAI API")
	}

	var results RecipeMetaResponseSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return RecipeMetaResponseSchema{}, errors.Wrap(err, "failed to unmarshal recipe response")
	}

	return results, nil
}
