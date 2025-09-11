package ai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/metrics"
	"github.com/openai/openai-go/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var GroceryListResponseJSONSchema = GenerateSchema[GroceryListResponseSchema]()

type GroceryListResponseSchema struct {
	GroceryListAggregatedIngredients []GroceryListIngredientSchema `json:"groceries"`
}

type GroceryListIngredientSchema struct {
	ID       int     `json:"id" jsonschema:"-"`
	AlikeID  int     `json:"alikeId" jsonschema_description:"Unique identifier for items that are the same type of ingredient. Items with same alikeId should have identical name, unit, and category. For example, all 'eggs' items should have the same alikeId regardless of which recipe they come from."`
	RecipeID string  `json:"recipeId" jsonschema_description:"The original recipeId that is passed in. DO NOT MODIFY OR MOVE THIS VALUE - KEEP IT EXACTLY THE SAME."`
	Name     string  `json:"name" jsonschema_description:"Generic, shopper-friendly name for the ingredient (e.g., 'eggs', 'tomatoes', 'ground beef'). Keep it lowercased."`
	Quantity float64 `json:"quantity" jsonschema_description:"Individual quantity for this specific ingredient instance. Do NOT aggregate across multiple ingredients."`
	Unit     string  `json:"unit" jsonschema_description:"Sensible, shopper-friendly unit for the ingredient (e.g., 'lbs', 'head', 'dozen'). For packaged items, specify the unit first, then the size in parentheses, like 'cans (15 oz)'."`
	Category string  `json:"category" jsonschema_description:"enum=produce,enum=dairy,enum=meat,enum=fish,enum=snacks,enum=canned_goods,enum=breads_and_bakery,enum=dry_and_baking_goods,enum=frozen"`
	IsMarked bool    `json:"isMarked" jsonschema:"-"`
}

func (a *AIClient) GeneratePlanGroceryList(ctx context.Context, planIngredientsMarkdown string) (GroceryListResponseSchema, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "grocerylist",
		Description: openai.String("A JSON object representing a grocery list"),
		Schema:      GroceryListResponseJSONSchema,
		Strict:      openai.Bool(true),
	}

	start := time.Now()
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       openai.ChatModelGPT4oMini,
		MaxTokens:   openai.Int(10000),
		Temperature: openai.Float(0.1),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(`
				CRITICAL: Transform each ingredient entry into exactly one grocery item. PRESERVE INDIVIDUAL QUANTITIES EXACTLY - NO AGGREGATION.

				INPUT-OUTPUT COUNT RULE: If you receive N ingredients, return EXACTLY N items.
				Example: 104 ingredients in → 104 groceries out. NO EXCEPTIONS.

				QUANTITY PRESERVATION RULE: Keep each individual quantity EXACTLY as provided in the input.
				- Input: 0.5 lbs ground beef → Output: 0.5 lbs ground beef (quantity preserved exactly)
				- Input: 0.5 lbs ground beef → Output: 0.5 lbs ground beef (quantity preserved exactly)
				- Input: 0.5 lbs ground beef → Output: 0.5 lbs ground beef (quantity preserved exactly)
				These should result in THREE separate items, each with quantity 0.5, NOT one item with quantity 1.5.

				UNIT STANDARDIZATION: Convert units to consistent, shopper-friendly formats, but DO NOT change quantities:
				- All weight measurements should use "lbs" when possible
				- All volume measurements should use standard cooking units (cups, tablespoons, etc.)
				- For items sold by count, use appropriate count units (whole, dozen, etc.)
				- For packaged items, use package-based units like "cans (15 oz)" or "boxes"

				FOR EACH INGREDIENT ENTRY, OUTPUT ONE GROCERY ITEM:

				1. ID: Sequential starting from 0
				2. ALIKE_ID: Assign the same alikeId ONLY for ingredients that are functionally identical for shopping purposes:
				   - "ground beef", "lean ground beef", "80/20 ground beef" → same alikeId (same shopping item)
				   - "ground beef" and "ground turkey" → different alikeIds (different products)
				   - "eggs", "large eggs", "chicken eggs" → same alikeId (same shopping item)
				   - "eggs" and "egg whites" → different alikeIds (different products)
				   - "yellow onions", "onions", "cooking onions" → same alikeId (same shopping item)
				   - "onions" and "green onions" → different alikeIds (different products)
				3. RECIPE_ID: Copy exactly from input - NEVER modify this field
				4. NAME: Use the most standard shopping name (e.g., "ground beef", "eggs", "onions")
				5. QUANTITY: Copy exactly from input - DO NOT aggregate or combine quantities
				6. UNIT: Standardize to shopper-friendly units (lbs, cups, dozen, whole, etc.)
				7. CATEGORY: produce, dairy, meat, fish, snacks, canned_goods, breads_and_bakery, dry_and_baking_goods, frozen

				EXAMPLES OF CORRECT BEHAVIOR:
				Input table with 3 rows of ground beef (0.5 lbs each):
				| recipe1 | 123 | ground beef | 0.5 | lbs | meat |
				| recipe2 | 124 | lean ground beef | 0.5 | pounds | meat |
				| recipe3 | 125 | 80/20 ground beef | 0.5 | lb | meat |

				Correct Output: THREE separate items:
				1. {"id": 0, "alikeId": 1, "recipeId": "recipe1", "name": "ground beef", "quantity": 0.5, "unit": "lbs", "category": "meat"}
				2. {"id": 1, "alikeId": 1, "recipeId": "recipe2", "name": "ground beef", "quantity": 0.5, "unit": "lbs", "category": "meat"}
				3. {"id": 2, "alikeId": 1, "recipeId": "recipe3", "name": "ground beef", "quantity": 0.5, "unit": "lbs", "category": "meat"}

				WRONG Output: One aggregated item (DO NOT DO THIS):
				{"id": 0, "alikeId": 1, "recipeId": "recipe1", "name": "ground beef", "quantity": 1.5, "unit": "lbs", "category": "meat"}

				Input data:
			` + planIngredientsMarkdown),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})
	if err != nil {
		return GroceryListResponseSchema{}, errors.Wrap(err, "unexpected error from OpenAI API")
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

	if len(chat.Choices) == 0 {
		return GroceryListResponseSchema{}, errors.New("no choices returned from OpenAI API")
	}

	var aiError error
	var results GroceryListResponseSchema
	for _, choice := range chat.Choices {
		err := json.Unmarshal([]byte(choice.Message.Content), &results)
		if err == nil {
			return results, nil
		}

		aiError = err
	}

	return GroceryListResponseSchema{}, aiError
}
