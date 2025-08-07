package folders

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type CreateFolderInput struct {
	FolderID  string
	KitchenID string
	Name      string
	Cover     null.String
}

func CreateFolder(ctx context.Context, store Store, input CreateFolderInput) (Folder, error) {
	_, err := store.ExecContext(ctx, `
		INSERT INTO folders (folder_id, kitchen_id, folder_name, cover, created_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, input.FolderID, input.KitchenID, input.Name, input.Cover)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") && strings.Contains(err.Error(), "key 'folders.folder_name'") {
			return Folder{}, ErrDuplicateFolderName
		}
		return Folder{}, err
	}

	return GetFolderByID(ctx, store, input.FolderID)
}

type UpdateFolderInput struct {
	FolderID string
	Name     string
	Cover    null.String
}

func UpdateFolder(ctx context.Context, store Store, input UpdateFolderInput) (Folder, error) {
	_, err := store.ExecContext(ctx, `
		UPDATE folders
		SET
			folder_name = ?,
			cover = ?
		WHERE
			folder_id = ?
	`, input.Name, input.Cover, input.FolderID)

	if err != nil {
		if err == sql.ErrNoRows {
			return Folder{}, ErrFolderNotFound
		}
		return Folder{}, err
	}

	return GetFolderByID(ctx, store, input.FolderID)
}

func ListFoldersByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Folder, error) {
	folders := make([]Folder, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT * FROM folders WHERE kitchen_id = ? ORDER BY created_at DESC
	`, kitchenID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var folder Folder
		if err := rows.StructScan(&folder); err != nil {
			return folders, err
		}
		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return folders, err
	}

	return folders, nil
}

func GetFolderByID(ctx context.Context, store Store, folderID string) (Folder, error) {
	var folder Folder
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM folders WHERE folder_id = ?
	`, folderID).StructScan(&folder)

	if err != nil {
		if err == sql.ErrNoRows {
			return Folder{}, ErrFolderNotFound
		}
		return Folder{}, err
	}

	return folder, nil
}

func DeleteFolderByID(ctx context.Context, store Store, folderID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM folders WHERE folder_id = ?
	`, folderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrFolderNotFound
		}
		return err
	}
	return nil
}
