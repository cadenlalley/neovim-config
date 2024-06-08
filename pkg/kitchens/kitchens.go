package kitchens

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type CreateKitchenInput struct {
	AccountID   string
	KitchenName string
	Bio         string
	Handle      string
	Avatar      string
	Cover       string
	Public      bool
}

func CreateKitchen(ctx context.Context, store Store, input CreateKitchenInput) (Kitchen, error) {
	kitchenID := CreateKitchenID()

	// Handle nullable values.
	bio := null.NewString(input.Bio, input.Bio != "")
	avatar := null.NewString(input.Avatar, input.Avatar != "")
	cover := null.NewString(input.Cover, input.Cover != "")

	_, err := store.ExecContext(ctx, `
		INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle, avatar, cover, public, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);
	`, kitchenID, input.AccountID, input.KitchenName, bio, input.Handle, avatar, cover, input.Public)

	if err != nil {
		return Kitchen{}, err
	}

	kitchen, err := GetKitchenByID(ctx, store, kitchenID)
	if err != nil {
		return Kitchen{}, err
	}

	return kitchen, nil
}

func GetKitchenByID(ctx context.Context, store Store, kitchenID string) (Kitchen, error) {
	var kitchen Kitchen
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM kitchens WHERE kitchen_id = ?;
	`, kitchenID).StructScan(&kitchen)

	if err != nil && err != sql.ErrNoRows {
		return Kitchen{}, err
	}

	return kitchen, nil
}
