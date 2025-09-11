package plans

import (
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Plan struct {
	ID             string       `json:"id" db:"meal_plan_id"`
	OwnerAccountID string       `json:"ownerAccountId" db:"account_id"`
	StartDate      time.Time    `json:"startDate" db:"start_date"`
	EndDate        time.Time    `json:"endDate" db:"end_date"`
	CreatedAt      time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time    `json:"updatedAt" db:"updated_at"`
	Recipes        []PlanRecipe `json:"recipes"`
}

type PlanRecipe struct {
	RecipeID    string `json:"recipeId" db:"recipe_id"`
	Day         int    `json:"day" db:"day"`
	ServingSize int    `json:"servingSize" db:"serving_size"`
}

type FullPlanRecipe struct {
	MealPlanRecipeID int         `json:"mealPlanRecipeId" db:"meal_plan_recipe_id"`
	RecipeID         string      `json:"recipeId" db:"recipe_id"`
	Day              int         `json:"day" db:"day_number"`
	PlanStartDate    time.Time   `json:"planStartDate" db:"start_date"`
	PlanEndDate      time.Time   `json:"planEndDate" db:"end_date"`
	PlannedDate      time.Time   `json:"plannedDate"`
	ServingSize      int         `json:"servingSize" db:"serving_size"`
	KitchenID        string      `json:"kitchenId" db:"kitchen_id"`
	Name             string      `json:"name" db:"recipe_name" validate:"required"`
	Summary          null.String `json:"summary" db:"summary"`
	PrepTime         *int        `json:"prepTime" db:"prep_time" validate:"required"`
	CookTime         *int        `json:"cookTime" db:"cook_time" validate:"required"`
	Servings         *int        `json:"servings" db:"servings" validate:"required"`
	Difficulty       int         `json:"difficulty" db:"difficulty"`
	Course           null.String `json:"course" db:"course"`
	Class            null.String `json:"class" db:"class"`
	Cuisine          null.String `json:"cuisine" db:"cuisine"`
	Cover            null.String `json:"cover" db:"cover"`
	Source           null.String `json:"source" db:"source"`
	CreatedAt        time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt        null.Time   `json:"deletedAt" db:"deleted_at"`
}

type GroceryListItem struct {
	ID       int     `json:"id" db:"item_id"`
	AlikeID  int     `json:"alikeId" db:"alike_id"`
	RecipeID string  `json:"recipeId" db:"recipe_id"`
	Name     string  `json:"name" db:"name"`
	Quantity float64 `json:"quantity" db:"quantity"`
	Unit     string  `json:"unit" db:"unit"`
	Category string  `json:"category" db:"category"`
	IsMarked bool    `json:"isMarked" db:"marked"`
}

type GroceryList struct {
	GroceryItems []GroceryListItem `json:"groceries"`
}

func CreatePlanID() string {
	return "pln_" + ksuid.New().String()
}
