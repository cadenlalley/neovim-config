package tags

import (
	"context"
	"database/sql"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

func CreateTags(ctx context.Context, store Store, inputs []Tag) ([]Tag, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	builder := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("tags").
		Columns("tag_type", "tag_value")

	for _, input := range inputs {
		t := strings.ToLower(input.Type)
		v := strings.ToLower(input.Value)

		// Skip invalid tag types.
		if err := ValidateTagType(t); err != nil {
			log.Error().Err(err).Str("tag_type", t).Msg("invalid tag type")
			continue
		}

		builder = builder.Values(t, v)
	}

	builder = builder.Suffix(`
		ON DUPLICATE KEY UPDATE
			tag_type = VALUES(tag_type),
			tag_value = VALUES(tag_value)
	`)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = store.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	// Get generated tags
	tags, err := GetTagsByTypeAndValue(ctx, store, inputs)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func GetTagsByTypeAndValue(ctx context.Context, store Store, inputs []Tag) ([]Tag, error) {
	builder := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Question).
		Select("tag_id", "tag_type", "tag_value").
		From("tags")

	orConditions := sq.Or{}
	for _, tag := range inputs {
		orConditions = append(orConditions, sq.Eq{
			"tag_type":  tag.Type,
			"tag_value": tag.Value,
		})
	}

	builder = builder.Where(orConditions)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.TagID, &tag.Type, &tag.Value)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
