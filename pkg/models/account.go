package models

import (
	"fmt"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Account struct {
	AccountID string    `json:"accountId" db:"account_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	CreatedAt null.Time `json:"createdAt" db:"created_at"`
	DeletedAt null.Time `json:"deletedAt" db:"deleted_at"`
}

func (a *Account) Exists() bool {
	return a.AccountID != ""
}

func CreateAccountID() string {
	return fmt.Sprintf("acc_%s", ksuid.New().String())
}
