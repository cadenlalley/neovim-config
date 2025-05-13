package recipes

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func CreateRecipeTags(ctx context.Context, store Store, recipeID string, tagIDs []int) error {
	builder := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("recipe_tags").
		Columns("recipe_id", "tag_id")

	for _, tagID := range tagIDs {
		builder = builder.Values(recipeID, tagID)
	}

	builder = builder.Suffix(`
		ON DUPLICATE KEY UPDATE
			recipe_id = VALUES(recipe_id),
			tag_id = VALUES(tag_id)
	`)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = store.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
