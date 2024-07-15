package kitchens

import (
	"errors"
)

var (
	ErrKitchenNotFound = errors.New("kitchen not found")
	ErrDuplicateHandle = errors.New("duplicate entry for handle")
)
