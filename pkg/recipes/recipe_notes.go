package recipes

import (
	"context"
)

type CreateRecipeNotesInput struct {
	RecipeID string
	StepID   int
	Note     string
}

func CreateRecipeNotes(ctx context.Context, store Store, input CreateRecipeNotesInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_notes (
			recipe_id,
			step_id,
			note
		) VALUES (?, ?, ?)
	`, input.RecipeID, input.StepID, input.Note)

	if err != nil {
		return nil
	}

	return nil
}

func DeleteRecipeNotesByRecipeID(ctx context.Context, store Store, recipeID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipe_notes WHERE recipe_id = ?;
	`, recipeID)

	if err != nil {
		return err
	}

	return nil
}

func GetRecipeNotesByRecipeID(ctx context.Context, store Store, recipeID string) ([]RecipeNote, error) {
	notes := make([]RecipeNote, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipe_notes WHERE recipe_id = ?
	`, recipeID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var note RecipeNote
		if err := rows.StructScan(&note); err != nil {
			return notes, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return notes, err
	}

	return notes, nil
}
