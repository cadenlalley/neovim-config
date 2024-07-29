package recipes

import (
	"context"
)

func GetRecipeImagesByRecipeID(ctx context.Context, store Store, recipeID string) ([]RecipeImage, error) {
	images := make([]RecipeImage, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM recipe_images WHERE recipe_id = ?
	`, recipeID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var image RecipeImage
		if err := rows.StructScan(&image); err != nil {
			return images, err
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return images, err
	}

	return images, nil
}
