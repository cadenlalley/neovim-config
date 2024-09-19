package recipes

import (
	"context"
)

type CreateRecipeStepInput struct {
	RecipeID    string
	StepID      int
	Instruction string
	Group       string
}

func CreateRecipeSteps(ctx context.Context, store Store, input CreateRecipeStepInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_steps (
			recipe_id,
			step_id,
			instruction,
			group_name
		) VALUES (?, ?, ?, ?)
	`, input.RecipeID, input.StepID, input.Instruction, input.Group)

	if err != nil {
		return err
	}

	return nil
}

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
