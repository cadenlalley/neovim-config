package sdk

import "github.com/kitchens-io/kitchens-api/pkg/followers"

// List Followers
func (c *client) ListKitchenFollowers(kitchenID string) ([]followers.Follower, error) {
	var followers []followers.Follower
	err := c.get("/v1/kitchen/"+kitchenID+"/followers", &followers)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// List Following
func (c *client) ListKitchenFollowing(kitchenID string) ([]followers.Follower, error) {
	var followed []followers.Follower
	err := c.get("/v1/kitchen/"+kitchenID+"/followed", &followed)
	if err != nil {
		return nil, err
	}
	return followed, nil
}
