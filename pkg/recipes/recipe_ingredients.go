package recipes

import (
	"context"

	"gopkg.in/guregu/null.v4"
)

type CreateRecipeIngredientInput struct {
	RecipeID     string
	IngredientID int
	Name         string
	Quantity     float64
	Unit         null.String
	Group        string
}

func CreateRecipeIngredients(ctx context.Context, store Store, input CreateRecipeIngredientInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_ingredients (
			recipe_id,
			ingredient_id,
			ingredient_name,
			quantity,
			unit,
			group_name
		) VALUES (?, ?, ?, ?, ?, ?)
	`, input.RecipeID, input.IngredientID, input.Name, input.Quantity, input.Unit, input.Group)

	if err != nil {
		return err
	}

	return nil
}

func DeleteRecipeIngredientsByRecipeID(ctx context.Context, store Store, recipeID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipe_ingredients WHERE recipe_id = ?;
	`, recipeID)

	if err != nil {
		return err
	}

	return nil
}

func DeleteRecipeIngredientsByKitchenID(ctx context.Context, store Store, kitchenID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipe_ingredients WHERE recipe_id IN (SELECT recipe_id FROM recipes WHERE kitchen_id = ?)
	`, kitchenID)
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
