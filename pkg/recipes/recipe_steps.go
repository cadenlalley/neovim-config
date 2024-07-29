package recipes

import (
	"context"
)

func GetRecipeStepsByRecipeID(ctx context.Context, store Store, recipeID string) ([]RecipeStep, error) {
	steps := make([]RecipeStep, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipe_steps WHERE recipe_id = ?
	`, recipeID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var step RecipeStep
		if err := rows.StructScan(&step); err != nil {
			return steps, err
		}
		steps = append(steps, step)
	}

	if err := rows.Err(); err != nil {
		return steps, err
	}

	return steps, nil
}
