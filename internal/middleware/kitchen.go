package middleware

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/auth"
	"github.com/kitchens-io/kitchens-api/pkg/folders"
	"github.com/kitchens-io/kitchens-api/pkg/kitchens"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
)

const KitchenContextKey = "kitchen"

type KitchenAuthorizer struct {
	db *sqlx.DB
}

func NewKitchenAuthorizer(db *sqlx.DB) *KitchenAuthorizer {
	return &KitchenAuthorizer{
		db: db,
	}
}

// Validate that the user has permission to this kitchen.
func (a *KitchenAuthorizer) ValidateWriter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID := c.Get(auth.UserIDContextKey).(string)
		kitchenID := c.Param("kitchen_id")

		// Lookup the user record for the provided JWT.
		account, err := accounts.GetAccountByUserID(ctx, a.db, userID)
		if err != nil {
			if err == accounts.ErrAccountNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "could not get account by user ID").SetInternal(err)
		}

		kitchen, err := kitchens.GetKitchenByID(ctx, a.db, kitchenID)
		if err != nil {
			if err == kitchens.ErrKitchenNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "could not get kitchen by ID").SetInternal(err)
		}

		// Validate that the user has permissions to be modifying this kitchen.
		if account.AccountID != kitchen.AccountID {
			err := fmt.Errorf("account '%s' attempted to create in kitchen '%s' without authorization", account.AccountID, kitchen.KitchenID)
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}

		// Store the Kitchen in the context.
		c.Set(KitchenContextKey, kitchen)

		return next(c)
	}
}

func (a *KitchenAuthorizer) ValidateFolderWriter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		folderID := c.Param("folder_id")

		kitchen := c.Get(KitchenContextKey).(kitchens.Kitchen)

		// Lookup folder resource.
		folder, err := folders.GetFolderByID(ctx, a.db, folderID)
		if err != nil {
			if err == folders.ErrFolderNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "could not get folder by ID").SetInternal(err)
		}

		// Validate that the user has permissions to be modifying this folder.
		if kitchen.KitchenID != folder.KitchenID {
			err := fmt.Errorf("kitchen '%s' attempted to modify a folder in kitchen '%s' without authorization", kitchen.KitchenID, folder.KitchenID)
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}

		return next(c)
	}
}

func (a *KitchenAuthorizer) ValidateRecipeWriter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		recipeID := c.Param("recipe_id")

		kitchen := c.Get(KitchenContextKey).(kitchens.Kitchen)

		// Lookup recipe Resource.
		recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
		if err != nil {
			if err == recipes.ErrRecipeNotFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipe by ID").SetInternal(err)
		}

		// Validate that the user has permissions to be modifying this recipe.
		if kitchen.KitchenID != recipe.KitchenID {
			err := fmt.Errorf("kitchen '%s' attempted to modify a recipe in kitchen '%s' without authorization", kitchen.KitchenID, recipe.KitchenID)
			return echo.NewHTTPError(http.StatusForbidden).SetInternal(err)
		}

		return next(c)
	}
}
