package recipes

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

// This regex matches any character that is NOT a-z, A-Z, 0-9, or dash (-)
var recipeShareURL = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

type Recipe struct {
	RecipeID   string      `json:"recipeId" db:"recipe_id"`
	KitchenID  string      `json:"kitchenId" db:"kitchen_id"`
	Name       string      `json:"name" db:"recipe_name" validate:"required"`
	Summary    null.String `json:"summary" db:"summary"`
	PrepTime   *int        `json:"prepTime" db:"prep_time" validate:"required"`
	CookTime   *int        `json:"cookTime" db:"cook_time" validate:"required"`
	Servings   *int        `json:"servings" db:"servings" validate:"required"`
	Difficulty int         `json:"difficulty" db:"difficulty"`
	Course     null.String `json:"course" db:"course"`
	Class      null.String `json:"class" db:"class"`
	Cuisine    null.String `json:"cuisine" db:"cuisine"`
	Cover      null.String `json:"cover" db:"cover"`
	Source     null.String `json:"source" db:"source"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt  null.Time   `json:"deletedAt" db:"deleted_at"`

	// Reviews
	ReviewCount  int     `json:"reviewCount" db:"review_count"`
	ReviewRating float64 `json:"reviewRating" db:"review_rating"`

	// Attached for full recipe
	SourceDomain null.String        `json:"sourceDomain" db:"-"`
	Ingredients  []RecipeIngredient `json:"ingredients" db:"-" validate:"required,dive"`
	Steps        []RecipeStep       `json:"steps" db:"-" validate:"required,dive"`
	ShareURL     string             `json:"shareUrl" db:"-"`
}

// Model validation not handled by the validator
func (r *Recipe) Validate() error {
	// TODO: Add validation for lower bound once UI implements.
	if r.Difficulty > 5 {
		return fmt.Errorf("field 'difficulty' must be between 1 and 5")
	}

	if len(r.Ingredients) == 0 {
		return fmt.Errorf("missing items for field 'ingredients'")
	}
	if len(r.Steps) == 0 {
		return fmt.Errorf("missing items for field 'steps'")
	}

	for _, i := range r.Ingredients {
		// NOTE: Fix for the UI sending 'Unit' as the default value.
		if i.Quantity.Float64 == 0 && (i.Unit.Valid && i.Unit.String != "") {
			return fmt.Errorf("ingredient '%s': field 'quantity' required when providing value for 'unit'", i.Name)
		}
	}

	return nil
}

// Handle computed values
func (r *Recipe) ComputeValues() error {
	// Source Domain
	if r.Source.Valid {
		parsedURL, err := url.Parse(r.Source.String)
		if err != nil {
			return err
		}
		host := parsedURL.Hostname()
		r.SourceDomain = null.NewString(host, host != "")
	}

	// Share URL
	r.ShareURL = r.FormatShareURL()

	return nil
}

// Format Share URL
func (r *Recipe) FormatShareURL() string {
	recipeName := strings.ToLower(recipeShareURL.ReplaceAllString(r.Name, "-"))
	recipeName = strings.TrimPrefix(recipeName, "-")
	recipeName = strings.TrimSuffix(recipeName, "-")

	recipeID := strings.Split(r.RecipeID, "rcp_")[1]

	return fmt.Sprintf("/%s/%s", recipeID, url.PathEscape(recipeName))
}

// Create from Import
func (r *Recipe) Import(v json.RawMessage, includeGroup bool) error {
	var input Recipe
	err := json.Unmarshal(v, &input)
	if err != nil {
		return err
	}

	r.Name = input.Name
	r.Summary = input.Summary
	r.PrepTime = input.PrepTime
	r.CookTime = input.CookTime
	r.Servings = input.Servings
	r.Difficulty = input.Difficulty
	r.Course = input.Course
	r.Class = input.Class
	r.Cuisine = input.Cuisine
	r.Source = input.Source
	r.Ingredients = make([]RecipeIngredient, len(input.Ingredients))
	r.Steps = make([]RecipeStep, len(input.Steps))

	for i, ingredient := range input.Ingredients {
		// FIXME: Temporary fix for import URL failing to handle groups well.
		group := ParseNullString(ingredient.Group)
		if !includeGroup {
			group = null.NewString("", false)
		}

		quantity := ParseNullFloat(ingredient.Quantity)

		// TODO: Combine these into a function that handles 1/3 and 2/3
		// Handle 1/3 quantity
		if quantity.Float64 >= 0.33 && quantity.Float64 < 0.34 {
			quantity = null.NewFloat(0.333, true)
		}

		// Handle 2/3 quantity
		if quantity.Float64 >= 0.66 && quantity.Float64 <= 0.67 {
			quantity = null.NewFloat(0.667, true)
		}

		r.Ingredients[i] = RecipeIngredient{
			IngredientID: ingredient.IngredientID,
			Name:         ingredient.Name,
			Quantity:     quantity,
			Unit:         ParseNullString(ingredient.Unit),
			Group:        group,
		}
	}

	for i, step := range input.Steps {
		// FIXME: Temporary fix for import URL failing to handle groups well.
		group := ParseNullString(step.Group)
		if !includeGroup {
			group = null.NewString("", false)
		}

		r.Steps[i] = RecipeStep{
			StepID:      step.StepID,
			Instruction: step.Instruction,
			Group:       group,
			Note:        step.Note,
		}
	}

	return nil
}

func CreateRecipeID() string {
	return "rcp_" + ksuid.New().String()
}

type RecipeStep struct {
	RecipeID    string      `json:"-" db:"recipe_id"`
	StepID      int         `json:"stepId" db:"step_id" validate:"required"`
	Instruction string      `json:"instruction" db:"instruction" validate:"required"`
	Group       null.String `json:"group" db:"group_name"`

	// Attached for full step
	Images []string `json:"images" db:"-"`
	Note   string   `json:"note" db:"-"`
}

type RecipeNote struct {
	RecipeID string `db:"recipe_id"`
	StepID   int    `db:"step_id"`
	Note     string `db:"note"`
}

type RecipeImage struct {
	RecipeID string `db:"recipe_id"`
	StepID   int    `db:"step_id"`
	ImageURL string `db:"image_url"`
}

type RecipeIngredient struct {
	RecipeID     string      `json:"-" db:"recipe_id"`
	IngredientID int         `json:"ingredientId" db:"ingredient_id" validate:"required"`
	Name         string      `json:"name" db:"ingredient_name" validate:"required"`
	Quantity     null.Float  `json:"quantity" db:"quantity"`
	Unit         null.String `json:"unit" db:"unit"`
	Group        null.String `json:"group" db:"group_name"`
}

type Review struct {
	ReviewID          string      `json:"reviewId" db:"review_id"`
	RecipeID          string      `json:"recipeId" db:"recipe_id"`
	ReviewerKitchenID string      `json:"reviewerKitchenId" db:"reviewer_kitchen_id"`
	ReviewerName      string      `json:"reviewerName" db:"reviewer_name"`
	ReviewerAvatar    null.String `json:"reviewerAvatar" db:"reviewer_avatar"`
	Description       null.String `json:"description" db:"review_description"`
	Rating            float64     `json:"rating" db:"rating"`
	Media             null.String `json:"media" db:"media_path"`
	CreatedAt         time.Time   `json:"createdAt" db:"created_at"`

	// Attached for full review
	TotalLikes int  `json:"totalLikes" db:"total_likes"`
	Liked      bool `json:"liked" db:"liked"`
}

func CreateReviewID() string {
	return "rvw_" + ksuid.New().String()
}

type ReviewSummary struct {
	Total         int         `json:"total" db:"total"`
	Average       float64     `json:"average" db:"average"`
	Rating_1      int         `json:"-" db:"rating_1"`
	Rating_2      int         `json:"-" db:"rating_2"`
	Rating_3      int         `json:"-" db:"rating_3"`
	Rating_4      int         `json:"-" db:"rating_4"`
	Rating_5      int         `json:"-" db:"rating_5"`
	Ratings       map[int]int `json:"ratings"`
	Reviews       []Review    `json:"reviews"`
	KitchenReview Review      `json:"kitchenReview"`
}

type SearchResult struct {
	RecipeID     string      `json:"recipeId" db:"recipe_id"`
	KitchenID    string      `json:"kitchenId" db:"kitchen_id"`
	Name         string      `json:"name" db:"recipe_name"`
	Cover        null.String `json:"cover" db:"cover"`
	Source       null.String `json:"source" db:"source"`
	ReviewCount  int         `json:"reviewCount" db:"review_count"`
	ReviewRating float64     `json:"reviewRating" db:"review_rating"`

	NameScore    float64 `json:"-" db:"name_score"`
	SummaryScore float64 `json:"-" db:"summary_score"`
}
