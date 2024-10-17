package recipes

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Recipe struct {
	RecipeID  string      `json:"recipeId" db:"recipe_id"`
	KitchenID string      `json:"kitchenId" db:"kitchen_id"`
	Name      string      `json:"name" db:"recipe_name" validate:"required"`
	Summary   null.String `json:"summary" db:"summary"`
	PrepTime  *int        `json:"prepTime" db:"prep_time" validate:"required"`
	CookTime  *int        `json:"cookTime" db:"cook_time" validate:"required"`
	Servings  *int        `json:"servings" db:"servings" validate:"required"`
	Cover     null.String `json:"cover" db:"cover"`
	Source    null.String `json:"source" db:"source"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time   `json:"deletedAt" db:"deleted_at"`

	// Attached for full recipe
	Ingredients []RecipeIngredient `json:"ingredients" db:"-" validate:"required,dive"`
	Steps       []RecipeStep       `json:"steps" db:"-" validate:"required,dive"`
}

// Model validation not handled by the validator
func (r *Recipe) Validate() error {
	if len(r.Ingredients) == 0 {
		return fmt.Errorf("missing items for field 'ingredients'")
	}
	if len(r.Steps) == 0 {
		return fmt.Errorf("missing items for field 'steps'")
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
	Quantity     float64     `json:"quantity" db:"quantity" validate:"required"`
	Unit         null.String `json:"unit" db:"unit"`
	Group        null.String `json:"group" db:"group_name"`
}
