package sdk

import (
	"github.com/kitchens-io/kitchens-api/pkg/accounts"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/kitchens-io/kitchens-api/pkg/tags"
)

// Admin List Accounts
func (c *client) AdminListAccounts() ([]accounts.Account, error) {
	var accounts []accounts.Account
	err := c.get("/admin/accounts", &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Admin Create Recipe
func (c *client) AdminCreateRecipe(kitchenID string, input recipes.Recipe) (recipes.Recipe, error) {
	var recipe recipes.Recipe
	err := c.post("/admin/kitchen/"+kitchenID+"/recipes", input, &recipe)
	if err != nil {
		return recipes.Recipe{}, err
	}
	return recipe, nil
}

// Admin Add Recipes to Folder
func (c *client) AdminAddFolderRecipes(kitchenID, folderID string, recipeIDs []string) error {
	type request struct {
		RecipeIDs []string `json:"recipeIds"`
	}

	err := c.post("/admin/kitchen/"+kitchenID+"/folders/"+folderID+"/recipes/add", request{
		RecipeIDs: recipeIDs,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Admin Follow Kitchen
func (c *client) AdminFollowKitchen(kitchenID, followedKitchenID string) error {
	type request struct {
		KitchenID string `json:"kitchenId"`
	}

	err := c.post("/admin/kitchen/"+kitchenID+"/follow", request{
		KitchenID: followedKitchenID,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Admin Create Recipe Metadata
func (c *client) AdminCreateRecipeMetadata(recipeID string) ([]tags.Tag, error) {
	type request struct {
		RecipeID string `json:"recipeId"`
	}

	var res []tags.Tag
	err := c.post("/admin/recipes/metadata", request{
		RecipeID: recipeID,
	}, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
