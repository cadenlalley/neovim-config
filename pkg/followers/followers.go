package followers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type FollowKitchenInput struct {
	KitchenID         string
	FollowedKitchenID string
}

func FollowKitchen(ctx context.Context, store Store, input FollowKitchenInput) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO kitchen_followers (kitchen_id, followed_kitchen_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP);
	`, input.KitchenID, input.FollowedKitchenID)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return nil
		}
		return err
	}

	return nil
}

type UnfollowKitchenInput struct {
	KitchenID         string
	FollowedKitchenID string
}

func UnfollowKitchen(ctx context.Context, store Store, input UnfollowKitchenInput) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM kitchen_followers WHERE kitchen_id = ? AND followed_kitchen_id = ?;
	`, input.KitchenID, input.FollowedKitchenID)

	if err != nil {
		return err
	}
	return nil
}
