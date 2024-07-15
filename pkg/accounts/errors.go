package accounts

import "errors"

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrDuplicateAccount = errors.New("account for user already exists")
)
