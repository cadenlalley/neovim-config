package recipes

import (
	"context"
	"strings"
)

type SaveRecipeInput struct {
	KitchenID string
	RecipeID  string
}

func SaveRecipe(ctx context.Context, store Store, input SaveRecipeInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipes_saved (kitchen_id, recipe_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP);
	`, input.KitchenID, input.RecipeID)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return nil
		}
		return err
	}
	return nil
}

type RemoveRecipeInput struct {
	KitchenID string
	RecipeID  string
}

func RemoveRecipe(ctx context.Context, store Store, input RemoveRecipeInput) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipes_saved WHERE kitchen_id = ? AND recipe_id = ?;
	`, input.KitchenID, input.RecipeID)

	if err != nil {
		return err
	}
	return nil
}
