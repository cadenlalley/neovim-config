package ai

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kitchens-io/kitchens-api/internal/metrics"
	"github.com/openai/openai-go/v2"
	"github.com/pkg/errors"
)

var RecipeResponseJSONSchema = GenerateSchema[RecipeResponseSchema]()

type RecipeResponseSchema struct {
	Name        string                            `json:"name"`
	Summary     string                            `json:"summary"`
	PrepTime    int                               `json:"prepTime" jsonschema_description:"the prep time in minutes"`
	CookTime    int                               `json:"cookTime" jsonschema_description:"the cook time in minutes"`
	Servings    int                               `json:"servings" jsonschema_description:"the number of servings"`
	Difficulty  int                               `json:"difficulty" jsonschema:"enum=1,enum=2,enum=3,enum=4,enum=5" jsonschema_description:"the difficulty of the recipe, 1 is the easiest and 5 is the hardest"`
	Course      string                            `json:"course" jsonschema:"enum=breakfast,enum=brunch,enum=lunch,enum=dinner,enum=dessert,enum=supper"`
	Class       string                            `json:"class" jsonschema:"enum=main,enum=side,enum=snack,enum=beverage,enum=dessert,enum=dip,enum=soup,enum=appetizer"`
	Cuisine     string                            `json:"cuisine" jsonschema_description:"the cuisine of the recipe"`
	Ingredients []RecipeResponseIngredientsSchema `json:"ingredients"`
	Steps       []RecipeResponseStepsSchema       `json:"steps"`
}

type RecipeResponseIngredientsSchema struct {
	IngredientID int     `json:"ingredientId"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity" jsonschema_description:"the amount of the ingredient, if an amount doesn't make sense for the ingredient set it to 0"`
	Unit         string  `json:"unit" jsonschema:"enum=bag,enum=bottle,enum=box,enum=can,enum=clove,enum=cup,enum=dash,enum=drop,enum=gallon,enum=gram,enum=jar,enum=kilogram,enum=liter,enum=milliliter,enum=ounce,enum=packet,enum=piece,enum=pint,enum=pinch,enum=pound,enum=quart,enum=slice,enum=stick,enum=tbsp,enum=tsp,enum=n/a" jsonschema_description:"optional unit of measurement, if a unit doesn't make sense for the ingredient (like whole vegetables) set it to n/a"`
	// Prepration   string  `json:"prepration" jsonschema_description:"optional preparation method, ex: peeled, chopped, minced, etc."`
	Group string `json:"group" jsonschema_description:"ingredient group or 'n/a' if the content does not subdivide ingredients"`
}

type RecipeResponseStepsSchema struct {
	StepID        int    `json:"stepId"`
	Instruction   string `json:"instruction"`
	Note          string `json:"note" jsonschema_description:"optional note for the instruction step"`
	Group         string `json:"group" jsonschema_description:"step group or 'n/a' if the content does not subdivide steps"`
	IngredientIDs []int  `json:"ingredientIds" jsonschema_description:"optional list of ingredient IDs that are used in this step"`
}

func (a *AIClient) ExtractRecipeFromText(ctx context.Context, text string) (RecipeResponseSchema, ResponseMetrics, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "recipe",
		Description: openai.String("A JSON object representing a recipe"),
		Schema:      RecipeResponseJSONSchema,
		Strict:      openai.Bool(true),
	}

	start := time.Now()
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:               openai.ChatModelGPT4oMini2024_07_18,
		MaxCompletionTokens: openai.Int(1600),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(`
Your task is to extract the recipe from the provided markdown text with perfect accuracy. Follow these rules carefully:

1. INGREDIENT HANDLING
   • Preserve all ingredient names, quantities, and units exactly as written
   • Maintain original ingredient grouping (e.g., "For the sauce", "For the crust")
   • Use empty string "" (not null) for ungrouped ingredients
   • Keep ingredients in their original order

2. STEP INSTRUCTIONS
   • Preserve the exact wording and numbering of steps
   • Keep steps in their original sequence
   • Remove decorative formatting (bold, italics) unless it's part of the instruction

3. INGREDIENT ASSOCIATION
   Include an ingredient in a step's ingredientIds ONLY when:
   • It's being actively used or transformed in that step
   • It's being combined with other ingredients for the first time
   • It's being added to the dish in its raw/prepared form

   Do NOT include an ingredient when:
   • It's part of a pre-mixed component from a previous step
   • It's only mentioned in reference to equipment (e.g., "greased pan")
   • It's part of a cooking instruction without being actively used (e.g., "bake for 30 minutes")

4. NOTES & TIPS
   • Move alternative methods, variations, and non-essential tips to the 'note' field
   • Keep the main instruction focused on the core action
   • Preserve cooking times and temperatures in the main instruction

5. SPECIAL CASES
   • For "combine all dry ingredients" type instructions, include all relevant ingredients
   • When an ingredient is prepped in one step (e.g., "chopped onions"), reference the base ingredient
   • For multi-component recipes, maintain clear separation between components

Example of good ingredient association:
Step: "Chop chocolate into small pieces"
ingredientIds: [1]  // chocolate

Step: "Whisk together flour, sugar, and salt"
ingredientIds: [2,3,4]  // flour, sugar, salt

Step: "Add milk and eggs, mix until combined"
ingredientIds: [5,6]    // milk, eggs

Step: "Bake for 30 minutes"
ingredientIds: []       // No ingredients actively used

---

**Recipe Markdown Text:**

` + text),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	if chat.Choices == nil {
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.New("no choices returned from OpenAI API")
	}

	var results RecipeResponseSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.Wrap(err, "failed to unmarshal recipe response")
	}

	return results, ResponseMetrics{
		Model:            chat.Model,
		PromptTokens:     chat.Usage.PromptTokens,
		CompletionTokens: chat.Usage.CompletionTokens,
		Latency:          metrics.Elapsed(start),
	}, nil
}

func (a *AIClient) ExtractRecipeFromImageURLs(ctx context.Context, urls []string) (RecipeResponseSchema, ResponseMetrics, error) {
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

	start := time.Now()
	chat, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:               openai.ChatModelGPT4oMini2024_07_18,
		MaxCompletionTokens: openai.Int(1600),
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
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.Wrap(err, "unexpected error from OpenAI API")
	}

	if len(chat.Choices) == 0 {
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.New("no choices returned from OpenAI API")
	}

	var results RecipeResponseSchema
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &results)
	if err != nil {
		return RecipeResponseSchema{}, ResponseMetrics{}, errors.Wrap(err, "failed to unmarshal recipe response")
	}

	return results, ResponseMetrics{
		Model:            chat.Model,
		PromptTokens:     chat.Usage.PromptTokens,
		CompletionTokens: chat.Usage.CompletionTokens,
		Latency:          metrics.Elapsed(start),
	}, nil
}
