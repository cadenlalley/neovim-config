package sdk

import "github.com/kitchens-io/kitchens-api/pkg/recipes"

// List Kitchen Recipes
func (c *client) ListKitchenRecipes(kitchenID string) ([]recipes.Recipe, error) {
	var recipes []recipes.Recipe
	err := c.get("/v1/kitchen/"+kitchenID+"/recipes", &recipes)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

// Get Kitchen Recipe
func (c *client) GetKitchenRecipe(kitchenID, recipeID string) (recipes.Recipe, error) {
	var recipe recipes.Recipe
	err := c.get("/v1/kitchen/"+kitchenID+"/recipes/"+recipeID, &recipe)
	if err != nil {
		return recipes.Recipe{}, err
	}
	return recipe, nil
}
