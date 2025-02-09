package followers

import "gopkg.in/guregu/null.v4"

type Follower struct {
	KitchenID string      `json:"kitchenId" db:"kitchen_id"`
	Name      string      `json:"name" db:"kitchen_name"`
	Handle    string      `json:"handle" db:"handle"`
	Avatar    null.String `json:"avatar" db:"avatar"`
}
