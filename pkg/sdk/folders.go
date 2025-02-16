package sdk

import "github.com/kitchens-io/kitchens-api/pkg/folders"

// List Folders
func (c *client) ListKitchenFolders(kitchenID string) ([]folders.Folder, error) {
	var folders []folders.Folder
	err := c.get("/v1/kitchen/"+kitchenID+"/folders", &folders)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

// Get Folder
func (c *client) GetKitchenFolder(kitchenID, folderID string) (folders.Folder, error) {
	var folder folders.Folder
	err := c.get("/v1/kitchen/"+kitchenID+"/folders/"+folderID, &folder)
	if err != nil {
		return folders.Folder{}, err
	}
	return folder, nil
}
