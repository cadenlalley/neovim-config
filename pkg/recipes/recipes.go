package recipes

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type CreateRecipeInput struct {
	RecipeID  string
	KitchenID string
	Name      string
	Summary   string
	PrepTime  int
	CookTime  int
	Servings  int
	Cover     string
	Source    string
}

func CreateRecipe(ctx context.Context, store Store, input CreateRecipeInput) (Recipe, error) {
	// Handle nullable values
	summary := null.NewString(input.Summary, input.Summary != "")
	cover := null.NewString(input.Cover, input.Cover != "")
	source := null.NewString(input.Source, input.Source != "")

	_, err := store.ExecContext(ctx, `
		INSERT INTO recipes (
			recipe_id,
			kitchen_id,
			recipe_name,
			summary,
			prep_time,
			cook_time,
			servings,
			cover,
			source,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);
	`, input.RecipeID, input.KitchenID, input.Name, summary, input.PrepTime, input.CookTime, input.Servings, cover, source)

	if err != nil {
		return Recipe{}, err
	}

	return GetRecipeByID(ctx, store, input.RecipeID)
}

func GetRecipeByID(ctx context.Context, store Store, recipeID string) (Recipe, error) {
	var recipe Recipe
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM recipes WHERE recipe_id = ?;
	`, recipeID).StructScan(&recipe)

	if err != nil {
		if err == sql.ErrNoRows {
			return Recipe{}, ErrRecipeNotFound
		}
		return Recipe{}, err
	}

	return recipe, nil
}

func ListRecipesByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Recipe, error) {
	recipes := make([]Recipe, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipes WHERE kitchen_id = ?
	`, kitchenID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var recipe Recipe
		if err := rows.StructScan(&recipe); err != nil {
			return recipes, err
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return recipes, err
	}

	return recipes, nil
}
