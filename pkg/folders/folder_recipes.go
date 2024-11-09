package folders

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type CreateFolderRecipeInput struct {
	FolderID string
	RecipeID string
}

func CreateFolderRecipe(ctx context.Context, store Store, input CreateFolderRecipeInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO folder_recipes (folder_id, recipe_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP);
	`, input.FolderID, input.RecipeID)

	if err != nil {
		return err
	}

	return nil
}

func ListFolderRecipesByFolderID(ctx context.Context, store Store, folderID string) ([]FolderRecipe, error) {
	folderRecipes := make([]FolderRecipe, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT
			r.recipe_id as recipe_id,
			r.recipe_name as recipe_name,
			r.cover as cover,
			fr.created_at as created_at
		FROM folder_recipes fr
			LEFT JOIN recipes r ON fr.recipe_id = r.recipe_id
		WHERE fr.folder_id = ?
		ORDER BY fr.created_at DESC
	`, folderID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var folderRecipe FolderRecipe
		if err := rows.StructScan(&folderRecipe); err != nil {
			return folderRecipes, err
		}
		folderRecipes = append(folderRecipes, folderRecipe)
	}

	if err := rows.Err(); err != nil {
		return folderRecipes, err
	}

	return folderRecipes, nil
}

func DeleteFolderRecipesByIDs(ctx context.Context, store Store, folderID string, recipeIDs []string) error {
	sql, args, err := sq.Delete("folder_recipes").Where(sq.Eq{
		"folder_id": folderID,
		"recipe_id": recipeIDs,
	}).ToSql()
	if err != nil {
		return err
	}

	_, err = store.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
