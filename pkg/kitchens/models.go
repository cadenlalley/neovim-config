package kitchens

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Kitchen struct {
	KitchenID string      `json:"kitchenId" db:"kitchen_id"`
	AccountID string      `json:"accountId" db:"account_id"`
	Name      string      `json:"name" db:"kitchen_name"`
	Bio       null.String `json:"bio" db:"bio"`
	Handle    string      `json:"handle" db:"handle"`
	Avatar    null.String `json:"avatarPhoto" db:"avatar"`
	Cover     null.String `json:"coverPhoto" db:"cover"`
	Private   bool        `json:"private" db:"is_private"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt null.Time   `json:"deletedAt" db:"deleted_at"`
}

func CreateKitchenID() string {
	return fmt.Sprintf("ktc_%s", ksuid.New().String())
}
