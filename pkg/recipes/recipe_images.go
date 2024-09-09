package recipes

import (
	"context"
)

type CreateRecipeImagesInput struct {
	RecipeID string
	StepID   int
	ImageURL string
}

func CreateRecipeImages(ctx context.Context, store Store, input CreateRecipeImagesInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_images (
			recipe_id,
			step_id,
			image_url
		) VALUES (?, ?, ?)
	`, input.RecipeID, input.StepID, input.ImageURL)

	if err != nil {
		return nil
	}

	return nil
}

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
