package sdk

import "github.com/kitchens-io/kitchens-api/pkg/kitchens"

// Search Kitchens
func (c *client) SearchKitchens() ([]kitchens.Kitchen, error) {
	var kitchens []kitchens.Kitchen
	err := c.get("/v1/kitchens/search", &kitchens)
	if err != nil {
		return nil, err
	}
	return kitchens, nil
}

func (c *client) GetKitchen(kitchenID string) (kitchens.Kitchen, error) {
	var kitchen kitchens.Kitchen
	err := c.get("/v1/kitchens/"+kitchenID, &kitchen)
	if err != nil {
		return kitchens.Kitchen{}, err
	}
	return kitchen, nil
}
