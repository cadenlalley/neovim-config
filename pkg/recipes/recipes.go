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
	Summary   null.String
	PrepTime  int
	CookTime  int
	Servings  int
	Cover     null.String
	Source    null.String
}

func CreateRecipe(ctx context.Context, store Store, input CreateRecipeInput) (Recipe, error) {
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
	`, input.RecipeID, input.KitchenID, input.Name, input.Summary, input.PrepTime, input.CookTime, input.Servings, input.Cover, input.Source)

	if err != nil {
		return Recipe{}, err
	}

	return GetRecipeByID(ctx, store, input.RecipeID)
}

type UpdateRecipeInput struct {
	RecipeID string
	Name     string
	Summary  null.String
	PrepTime int
	CookTime int
	Servings int
	Cover    null.String
	Source   null.String
}

func UpdateRecipe(ctx context.Context, store Store, input UpdateRecipeInput) (Recipe, error) {
	_, err := store.ExecContext(ctx, `
		UPDATE recipes
		SET
			recipe_name = ?,
			summary = ?,
			prep_time = ?,
			cook_time = ?,
			servings = ?,
			cover = ?,
			source = ?
		WHERE
			recipe_id = ?;
	`, input.Name, input.Summary, input.PrepTime, input.CookTime, input.Servings, input.Cover, input.Source, input.RecipeID)
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

	// TODO: Revisit if this is the best way to handle this.
	err = recipe.ComputeValues()
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func ListRecipesByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Recipe, error) {
	recipes := make([]Recipe, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipes
		WHERE (kitchen_id = ? OR recipe_id IN (SELECT recipe_id FROM recipes_saved WHERE kitchen_id = ?))
			AND deleted_at IS NULL
		ORDER BY created_at;
	`, kitchenID, kitchenID)

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

func DeleteRecipeByID(ctx context.Context, store Store, recipeID string) error {
	_, err := store.ExecContext(ctx, `
		UPDATE recipes SET deleted_at = CURRENT_TIMESTAMP WHERE recipe_id = ?
	`, recipeID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRecipesByKitchenID(ctx context.Context, store Store, kitchenID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipes WHERE kitchen_id = ?
	`, kitchenID)
	if err != nil {
		return err
	}
	return nil
}
