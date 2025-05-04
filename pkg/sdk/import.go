package sdk

import (
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
)

// Import Recipe
func (c *client) ImportRecipeFromURL(url string) (recipes.Recipe, error) {
	type req struct {
		Source string `json:"source"`
	}

	var recipe recipes.Recipe
	err := c.post("/v1/import/url", req{
		Source: url,
	}, &recipe)
	if err != nil {
		return recipes.Recipe{}, err
	}

	return recipe, nil
}
