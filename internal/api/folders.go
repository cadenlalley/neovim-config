package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/override"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/folders"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type CreateKitchenFolderRequest struct {
	Name string `form:"name" validate:"required"`

	// The following are manually checked in the CreateFolder handler.
	// They cannot be bound automatically, and are optional.
	//
	// folderCoverFile
}

func (a *App) CreateKitchenFolder(c echo.Context) error {
	var input CreateKitchenFolderRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	// Create folder ID
	folderID := folders.CreateFolderID()

	prefix := media.GetFolderMediaPath(folderID)
	folderKey, err := a.handleFormFile(c, "folderCoverFile", prefix)
	if err != nil && err != http.ErrMissingFile {
		err = errors.Wrapf(err, "could not upload folderCoverFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload folder cover").SetInternal(err)
	}

	folder, err := folders.CreateFolder(ctx, a.db, folders.CreateFolderInput{
		FolderID:  folderID,
		KitchenID: kitchenID,
		Name:      input.Name,
		Cover:     null.NewString(folderKey, folderKey != ""),
	})
	if err != nil {
		if err == folders.ErrDuplicateFolderName {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create folder").SetInternal(err)
	}

	return c.JSON(http.StatusOK, folder)
}

type UpdateKitchenFolderRequest struct {
	Name string `form:"name"`

	// The following are manually checked in the UpdateFolder handler.
	// They cannot be bound automatically, and are optional.
	//
	// folderCoverFile
}

func (a *App) UpdateKitchenFolder(c echo.Context) error {
	var input UpdateKitchenFolderRequest
	err := web.ValidateRequest(c, web.ContentTypeMultipartFormData, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	folderID := c.Param("folder_id")

	folder, err := folders.GetFolderByID(ctx, a.db, folderID)
	if err != nil {
		if err == folders.ErrFolderNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "folder not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get folder").SetInternal(err)
	}

	// Handle the file uploads if they have been set.
	prefix := media.GetFolderMediaPath(folderID)
	folderKey, err := a.handleFormFile(c, "folderCoverFile", prefix)
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		err = errors.Wrapf(err, "could not upload folderCoverFile to prefix '%s'", prefix)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not upload folder cover").SetInternal(err)
	}

	// Since the folderFile is allowed to be null
	folderFile := folder.Cover.String
	if folderKey != "" {
		folderFile = folderKey
	}

	folder, err = folders.UpdateFolder(ctx, a.db, folders.UpdateFolderInput{
		FolderID: folderID,
		Name:     override.String(input.Name, folder.Name),
		Cover:    override.NullString(folderFile, folder.Cover),
	})
	if err != nil {
		if err == folders.ErrDuplicateFolderName {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update folder").SetInternal(err)
	}

	return c.JSON(http.StatusOK, folder)
}

func (a *App) GetKitchenFolder(c echo.Context) error {
	ctx := c.Request().Context()
	folderID := c.Param("folder_id")

	folder, err := folders.GetFolderByID(ctx, a.db, folderID)
	if err != nil {
		if err == folders.ErrFolderNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "folder not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get folder").SetInternal(err)
	}

	recipes, err := folders.ListFolderRecipesByFolderID(ctx, a.db, folderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not list folder recipes").SetInternal(err)
	}

	folder.Recipes = recipes

	return c.JSON(http.StatusOK, folder)
}

func (a *App) GetKitchenFolders(c echo.Context) error {
	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	folderList, err := folders.ListFoldersByKitchenID(ctx, a.db, kitchenID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not list folders").SetInternal(err)
	}

	return c.JSON(http.StatusOK, folderList)
}

func (a *App) DeleteKitchenFolder(c echo.Context) error {
	ctx := c.Request().Context()
	folderID := c.Param("folder_id")

	txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		err := folders.DeleteFolderRecipesByFolderID(ctx, tx, folderID)
		if err != nil {
			return errors.Wrap(err, "could not delete folder recipes")
		}

		err = folders.DeleteFolderByID(ctx, tx, folderID)
		if err != nil {
			return errors.Wrap(err, "could not delete folder")
		}

		return nil
	})
	if txErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete folder").SetInternal(txErr)
	}

	return c.NoContent(http.StatusNoContent)
}

type CreateKitchenFolderRecipeRequest struct {
	RecipeIDs []string `json:"recipeIds" validate:"required,min=1"`
}

func (a *App) CreateKitchenFolderRecipes(c echo.Context) error {
	var input CreateKitchenFolderRecipeRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	folderID := c.Param("folder_id")

	err = folders.CreateFolderRecipes(ctx, a.db, folderID, input.RecipeIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create folder recipes").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type DeleteKitchenFolderRecipesRequest struct {
	RecipeIDs []string `json:"recipeIds" validate:"required,min=1"`
}

func (a *App) DeleteKitchenFolderRecipes(c echo.Context) error {
	var input DeleteKitchenFolderRecipesRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	folderID := c.Param("folder_id")

	err = folders.DeleteFolderRecipesByIDs(ctx, a.db, folderID, input.RecipeIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete folder recipes").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
