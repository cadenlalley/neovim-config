package api

import (
	"net/http"

	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/followers"
	"github.com/labstack/echo/v4"
)

func (a *App) GetKitchenFollowers(c echo.Context) error {
	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	followersList, err := followers.ListFollowersByKitchenID(ctx, a.db, kitchenID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, followersList)
}

func (a *App) GetKitchensFollowing(c echo.Context) error {
	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	followedList, err := followers.ListFollowingByKitchenID(ctx, a.db, kitchenID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, followedList)
}

type FollowKitchenRequest struct {
	KitchenID string `json:"kitchenId" validate:"required"`
}

func (a *App) FollowKitchen(c echo.Context) error {
	var input FollowKitchenRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	err = followers.FollowKitchen(ctx, a.db, followers.FollowKitchenInput{
		KitchenID:         kitchenID,
		FollowedKitchenID: input.KitchenID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type UnfollowKitchenRequest struct {
	KitchenID string `json:"kitchenId" validate:"required"`
}

func (a *App) UnfollowKitchen(c echo.Context) error {
	var input UnfollowKitchenRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	kitchenID := c.Param("kitchen_id")

	err = followers.UnfollowKitchen(ctx, a.db, followers.UnfollowKitchenInput{
		KitchenID:         kitchenID,
		FollowedKitchenID: input.KitchenID,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
