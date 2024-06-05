package accounts

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/pkg/models"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type CreateAccountInput struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

func CreateAccount(ctx context.Context, store Store, input CreateAccountInput) (models.Account, error) {
	accountID := models.CreateAccountID()

	_, err := store.ExecContext(ctx, `
		INSERT INTO accounts (account_id, user_id, email, first_name, last_name, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, accountID, input.UserID, input.Email, input.FirstName, input.LastName)

	if err != nil {
		return models.Account{}, err
	}

	account, err := GetAccountByUserID(ctx, store, input.UserID)
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func GetAccountByID(ctx context.Context, store Store, accountID string) (models.Account, error) {
	var account models.Account
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM accounts WHERE account_id = ?
	`, accountID).StructScan(&account)

	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func GetAccountByUserID(ctx context.Context, store Store, userID string) (models.Account, error) {
	var account models.Account
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM accounts WHERE user_id = ?
	`, userID).StructScan(&account)

	if err != nil && err != sql.ErrNoRows {
		return models.Account{}, err
	}

	return account, nil
}
