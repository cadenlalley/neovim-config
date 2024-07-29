package kitchens

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type CreateKitchenInput struct {
	AccountID   string
	KitchenName string
	Bio         string
	Handle      string
	Avatar      string
	Cover       string
	Private     bool
}

func CreateKitchen(ctx context.Context, store Store, input CreateKitchenInput) (Kitchen, error) {
	kitchenID := CreateKitchenID()

	// Handle nullable values.
	bio := null.NewString(input.Bio, input.Bio != "")
	avatar := null.NewString(input.Avatar, input.Avatar != "")
	cover := null.NewString(input.Cover, input.Cover != "")

	_, err := store.ExecContext(ctx, `
		INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle, avatar, cover, is_private, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);
	`, kitchenID, input.AccountID, input.KitchenName, bio, input.Handle, avatar, cover, input.Private)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") && strings.Contains(err.Error(), "key 'kitchens.handle'") {
			return Kitchen{}, ErrDuplicateHandle
		}
		return Kitchen{}, err
	}

	return GetKitchenByID(ctx, store, kitchenID)
}

type UpdateKitchenInput struct {
	KitchenID string
	Name      string
	Bio       null.String
	Handle    string
	Avatar    null.String
	Cover     null.String
	Private   bool
}

func UpdateKitchen(ctx context.Context, store Store, input UpdateKitchenInput) (Kitchen, error) {
	_, err := store.ExecContext(ctx, `
		UPDATE kitchens
		SET
			kitchen_name = ?,
			bio = ?,
			handle = ?,
			avatar = ?,
			cover = ?,
			is_private = ?
		WHERE
			kitchen_id = ?
	`, input.Name, input.Bio, input.Handle, input.Avatar, input.Cover, input.Private, input.KitchenID)
	if err != nil {
		return Kitchen{}, err
	}

	return GetKitchenByID(ctx, store, input.KitchenID)
}

func GetKitchenByID(ctx context.Context, store Store, kitchenID string) (Kitchen, error) {
	var kitchen Kitchen
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM kitchens WHERE kitchen_id = ?;
	`, kitchenID).StructScan(&kitchen)

	if err != nil {
		if err == sql.ErrNoRows {
			return Kitchen{}, ErrKitchenNotFound
		}
		return Kitchen{}, err
	}

	return kitchen, nil
}

func ListKitchensByAccountID(ctx context.Context, store Store, accountID string) ([]Kitchen, error) {
	kitchens := make([]Kitchen, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM kitchens WHERE account_id = ?
	`, accountID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var kitchen Kitchen
		if err := rows.StructScan(&kitchen); err != nil {
			return kitchens, err
		}
		kitchens = append(kitchens, kitchen)
	}

	if err := rows.Err(); err != nil {
		return kitchens, err
	}

	return kitchens, nil
}
