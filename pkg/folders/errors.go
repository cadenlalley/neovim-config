package folders

import "errors"

var (
	ErrFolderNotFound      = errors.New("folder not found")
	ErrDuplicateFolderName = errors.New("duplicate entry for folder")
)
