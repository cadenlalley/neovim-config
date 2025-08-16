package sdk

import (
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
)

// Import Recipe
func (c *client) ImportRecipeFromURL(url string) (recipes.Recipe, *ai.ResponseMetrics, error) {
	type req struct {
		Source   string `json:"source"`
		Features struct {
			Groups bool `json:"groups"`
		} `json:"features"`
		Debug bool `json:"debug"`
	}

	type res struct {
		*recipes.Recipe
		Metrics *ai.ResponseMetrics `json:"metrics,omitempty"`
	}

	var response res
	err := c.post("/v1/import/url", req{
		Source: url,
		Debug:  true,
		Features: struct {
			Groups bool `json:"groups"`
		}{
			Groups: true,
		},
	}, &response)
	if err != nil {
		return recipes.Recipe{}, nil, err
	}

	return *response.Recipe, response.Metrics, nil
}
