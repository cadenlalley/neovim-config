package api

import (
	"net/http"
	"time"

	"github.com/kitchens-io/kitchens-api/pkg/plans"
	"github.com/labstack/echo/v4"
)

type CreatePlanRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (a *App) CreatePlan(c echo.Context) error {
	ctx := c.Request().Context()

	accountID := c.Param("account_id")

	var req CreatePlanRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	startDateParsed, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse start date").SetInternal(err)
	}

	endDateParsed, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse end date").SetInternal(err)
	}

	plan, err := plans.CreatePlan(ctx, a.db, plans.CreatePlanInput{
		ID:             plans.CreatePlanID(),
		OwnerAccountID: accountID,
		StartDate:      startDateParsed,
		EndDate:        endDateParsed,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create plan").SetInternal(err)
	}

	return c.JSON(http.StatusOK, plan)
}

type GetPlansByUserIDResponse struct {
	Plans []plans.Plan `json:"plans"`
}

func (a *App) GetPlansByUserID(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("account_id")

	plans, err := plans.ListPlansByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get plans by user ID").SetInternal(err)
	}

	return c.JSON(http.StatusOK, GetPlansByUserIDResponse{
		Plans: plans,
	})
}

func (a *App) GetPlanByAccountIDAndDateRange(c echo.Context) error {
	ctx := c.Request().Context()

	accountID := c.Param("account_id")
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse start date").SetInternal(err)
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse end date").SetInternal(err)
	}

	plan, err := plans.GetPlanByAccountIDAndDateRange(ctx, a.db, accountID, startDateParsed, endDateParsed)
	if err != nil {
		if err == plans.ErrPlanNotFound {
			return c.NoContent(http.StatusOK)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get plans").SetInternal(err)
	}

	recipes, err := plans.GetRecipesByPlanID(ctx, a.db, plan.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipes for plan").SetInternal(err)
	}

	plan.Recipes = recipes

	return c.JSON(http.StatusOK,
		plan,
	)
}

type GetPlanByIDResponse struct {
	Plan plans.Plan `json:"plan"`
}

func (a *App) GetPlanByID(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	plan, err := plans.GetPlanByID(ctx, a.db, planID)
	if err != nil {
		if err == plans.ErrPlanNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get plan by ID").SetInternal(err)
	}

	c.JSON(http.StatusOK, GetPlanByIDResponse{
		Plan: plan,
	})

	return nil
}

func (a *App) GetUserPlans(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("account_id")

	plans, err := plans.ListPlansByUserID(ctx, a.db, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get plans by user ID").SetInternal(err)
	}

	return c.JSON(http.StatusOK, plans)
}

type AddRecipesToPlanRequest struct {
	RecipeID    string `json:"recipeId"`
	Days        []int  `json:"days"`
	ServingSize int    `json:"servingSize"`
}

func (a *App) AddRecipesToPlan(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	var req AddRecipesToPlanRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.Days == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no days provided")
	}

	for _, day := range req.Days {
		err = plans.AddRecipesToPlan(ctx, a.db, planID, req.RecipeID, day, req.ServingSize)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not add recipe to plan").SetInternal(err)
		}
	}

	return c.NoContent(http.StatusOK)
}

type GetFullRecipesByPlanIDResponse struct {
	RecipeMap map[int][]plans.FullPlanRecipe `json:"recipeMap"`
}

func (a *App) GetFullRecipesByPlanID(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	recipes, err := plans.GetFullRecipesByPlanID(ctx, a.db, planID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get recipes for plan").SetInternal(err)
	}

	recipeMap := map[int][]plans.FullPlanRecipe{
		0: {},
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
	}

	var startDate time.Time
	for _, recipe := range recipes {
		startDate = recipe.PlanStartDate
		plannedDate := recipe.PlanStartDate.Add(time.Duration(recipe.Day) * 24 * time.Hour)
		recipe.PlannedDate = plannedDate

		recipeMap[recipe.Day] = append(recipeMap[recipe.Day], recipe)
	}

	for i := range 6 {
		for _, recipe := range recipeMap[i] {
			if recipe.RecipeID != "" {
				continue
			}

			recipe = plans.FullPlanRecipe{
				Day:           i,
				PlanStartDate: startDate,
				PlanEndDate:   startDate.AddDate(0, 0, 6),
				PlannedDate:   startDate.AddDate(0, 0, i),
			}
		}
	}

	return c.JSON(http.StatusOK, GetFullRecipesByPlanIDResponse{
		RecipeMap: recipeMap,
	})
}
