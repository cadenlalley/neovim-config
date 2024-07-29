package recipes

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Recipe struct {
	RecipeID  string      `json:"recipeId" db:"recipe_id"`
	KitchenID string      `json:"kitchenId" db:"kitchen_id"`
	Name      string      `json:"name" db:"recipe_name"`
	Summary   null.String `json:"summary" db:"summary"`
	PrepTime  int         `json:"prepTime" db:"prep_time"`
	CookTime  int         `json:"cookTime" db:"cook_time"`
	Servings  int         `json:"servings" db:"servings"`
	Cover     null.String `json:"cover" db:"cover"`
	Source    null.String `json:"source" db:"source"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time   `json:"deletedAt" db:"deleted_at"`

	// Attached for full recipe
	Steps []RecipeStep `json:"steps"`
}

type RecipeStep struct {
	RecipeID    string `json:"-" db:"recipe_id"`
	StepID      int    `json:"stepId" db:"step_id"`
	Instruction string `json:"instruction" db:"instruction"`

	// Attached for full step
	Images []string    `json:"images"`
	Notes  null.String `json:"notes"`
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

// type RecipeIngredient struct {
// 	RecipeID     string  `json:"recipeId" db:"recipe_id"`
// 	IngredientID int     `json:"ingredientId" db:"ingredient_id"`
// 	Name         string  `json:"name" db:"name"`
// 	Quantity     float64 `json:"quantity" db:"quantity"`
// 	Unit         string  `json:"unit" db:"unit"`
// }
