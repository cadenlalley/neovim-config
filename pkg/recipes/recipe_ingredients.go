package recipes

import (
	"context"
)

type CreateRecipeIngredientInput struct {
	RecipeID     string
	IngredientID int
	Name         string
	Quantity     float64
	Unit         string
}

func CreateRecipeIngredients(ctx context.Context, store Store, input CreateRecipeIngredientInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_ingredients (
			recipe_id,
			ingredient_id,
			ingredient_name,
			quantity,
			unit
		) VALUES (?, ?, ?, ?, ?)
	`, input.RecipeID, input.IngredientID, input.Name, input.Quantity, input.Unit)

	if err != nil {
		return err
	}

	return nil
}

func GetRecipeIngredientsByRecipeID(ctx context.Context, store Store, recipeID string) ([]RecipeIngredient, error) {
	ingredients := make([]RecipeIngredient, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipe_ingredients WHERE recipe_id = ?
	`, recipeID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ingredient RecipeIngredient
		if err := rows.StructScan(&ingredient); err != nil {
			return ingredients, err
		}
		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		return ingredients, err
	}

	return ingredients, nil
}
