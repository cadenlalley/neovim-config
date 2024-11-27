package folders

import (
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/guregu/null.v4"
)

type Folder struct {
	FolderID  string      `json:"folderId" db:"folder_id"`
	KitchenID string      `json:"kitchenId" db:"kitchen_id" validate:"required"`
	Name      string      `json:"name" db:"folder_name" validate:"required"`
	Cover     null.String `json:"cover" db:"cover"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`

	// Attached for full response
	Recipes []FolderRecipe `json:"recipes" db:"-"`
}

type FolderRecipe struct {
	RecipeID  string      `json:"recipeId" db:"recipe_id" validate:"required"`
	Name      string      `json:"name" db:"recipe_name"`
	Cover     null.String `json:"cover" db:"cover"`
	CreatedAt string      `json:"createdAt" db:"created_at"`
}

func CreateFolderID() string {
	return "fld_" + ksuid.New().String()
}
