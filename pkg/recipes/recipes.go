package recipes

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
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
