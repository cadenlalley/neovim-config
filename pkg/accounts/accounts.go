package accounts

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type CreateAccountInput struct {
	AccountID string
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

func CreateAccount(ctx context.Context, store Store, input CreateAccountInput) (Account, error) {
	_, err := store.ExecContext(ctx, `
		INSERT INTO accounts (account_id, user_id, email, first_name, last_name, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, input.AccountID, input.UserID, input.Email, input.FirstName, input.LastName)

	if err != nil {
		return Account{}, err
	}

	return GetAccountByID(ctx, store, input.AccountID)
}

type UpdateAccountInput struct {
	AccountID string
	FirstName string
	LastName  string
}

func UpdateAccount(ctx context.Context, store Store, input UpdateAccountInput) (Account, error) {
	_, err := store.ExecContext(ctx, `
		UPDATE accounts
		SET
			first_name = ?,
			last_name = ?
		WHERE
			account_id = ?
	`, input.FirstName, input.LastName, input.AccountID)
	if err != nil {
		return Account{}, err
	}

	return GetAccountByID(ctx, store, input.AccountID)
}

func ListAccounts(ctx context.Context, store Store) ([]Account, error) {
	accounts := make([]Account, 0)

	rows, err := store.QueryxContext(ctx, `SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account Account
		if err := rows.StructScan(&account); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccountByID(ctx context.Context, store Store, accountID string) (Account, error) {
	var account Account
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM accounts WHERE account_id = ?
	`, accountID).StructScan(&account)

	if err != nil {
		if err == sql.ErrNoRows {
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}

	return account, nil
}

func GetAccountByUserID(ctx context.Context, store Store, userID string) (Account, error) {
	var account Account
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM accounts WHERE user_id = ?
	`, userID).StructScan(&account)

	if err != nil {
		if err == sql.ErrNoRows {
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}

	return account, nil
}

func DeleteAccountByID(ctx context.Context, store Store, accountID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM accounts WHERE account_id = ?
	`, accountID)

	if err != nil {
		return err
	}

	return nil
}
