package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/web"
	"github.com/kitchens-io/kitchens-api/pkg/plans"
	recipestore "github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type CreatePlanRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (a *App) CreatePlan(c echo.Context) error {
	ctx := c.Request().Context()

	accountID := c.Param("account_id")

	var req CreatePlanRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &req)
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
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &req)
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

	err = plans.SetPlanGroceryListIsDirty(ctx, a.db, planID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not mark grocery list as dirty").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) RemoveRecipeFromPlan(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")
	recipeID := c.Param("recipe_id")

	err := plans.RemoveRecipeFromPlan(ctx, a.db, recipeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not remove recipe from plan").SetInternal(err)
	}

	err = plans.SetPlanGroceryListIsDirty(ctx, a.db, planID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not mark grocery list as dirty").SetInternal(err)
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

	recipeMap := make(map[int][]plans.FullPlanRecipe)

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

type aggregatedGroceryListResponse struct {
	GroceryList ai.GroceryListResponseSchema `json:"groceryList"`
}

func (a *App) GetGroceryListByPlanID(c echo.Context) error {
	// Use same timeout as OpenAI requests for grocery list generation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*240)
	defer cancel()

	planID := c.Param("id")

	// medium path - meal plan has been updated since last grocery list update
	isDirty, err := plans.GetPlanGroceryListIsDirty(ctx, a.db, planID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not check if meal plan has been updated since last grocery list generation").SetInternal(err)
	}

	if isDirty == true {
		txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
			err := plans.DeleteGroceryList(ctx, tx, planID)
			if err != nil {
				return errors.Wrap(err, "could not delete outdated grocery list")
			}

			aiGeneratedGroceryList, err := createGroceryList(ctx, planID, a)
			if err != nil {
				return errors.Wrap(err, "could not create grocery list")
			}

			err = plans.CreateGroceryList(ctx, tx, planID, aiGeneratedGroceryList)
			if err != nil {
				return errors.Wrap(err, "could not store updated grocery list")
			}

			err = plans.SetPlanGroceryListIsDirty(ctx, tx, planID, false)
			if err != nil {
				return errors.Wrap(err, "could not update grocery lsit dirty status")
			}

			return nil
		})

		if txErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not regenerate grocery list").SetInternal(err)
		}
	}

	// quick path - grocery list items are already in the db, do not need to make a request to open ai
	storedGroceryList, err := plans.GetGroceryListByMealPlanID(ctx, a.db, planID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting preexisting grocery list").SetInternal(err)
	}

	// slow path - no grocery list in db, need to have openai generate it
	if storedGroceryList == nil {
		// using the same model that the ai generates feels a bit hacky but might also make sense to keep api contracts the same without struct field swapping everywhere
		aiGeneratedGroceryList, err := createGroceryList(ctx, planID, a)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not create grocery list").SetInternal(err)
		}

		err = plans.CreateGroceryList(ctx, a.db, planID, aiGeneratedGroceryList)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store grocery list").SetInternal(err)
		}

		return c.JSON(http.StatusOK, aiGeneratedGroceryList)
	}

	categoryOrder, err := plans.GetCategoryOrderByPlanID(ctx, a.db, planID)
	if err != nil {
		return c.JSON(http.StatusOK, storedGroceryList)
	}

	storedGroceryList.CategoryOrder = categoryOrder

	return c.JSON(http.StatusOK, storedGroceryList)
}

func createGroceryList(ctx context.Context, planID string, app *App) (plans.GroceryList, error) {
	recipes, err := plans.GetRecipesByPlanID(ctx, app.db, planID)
	if err != nil {
		return plans.GroceryList{}, err
	}

	// Structure to include recipe ID with each ingredient
	type IngredientWithRecipeID struct {
		RecipeID     string  `json:"recipeId"`
		IngredientID int     `json:"ingredientId"`
		Name         string  `json:"name"`
		Quantity     float64 `json:"quantity"`
		Unit         string  `json:"unit"`
		Group        string  `json:"group"`
	}

	var ingredientsWithRecipeIDs []IngredientWithRecipeID
	for _, recipe := range recipes {
		originalRecipe, err := recipestore.GetRecipeByID(ctx, app.db, recipe.RecipeID)
		if err != nil {
			return plans.GroceryList{}, err
		}

		ingredients, err := recipestore.GetRecipeIngredientsByRecipeID(ctx, app.db, recipe.RecipeID)
		if err != nil {
			return plans.GroceryList{}, err
		}

		for _, ing := range ingredients {
			adjustedQuantity := ing.Quantity
			if originalRecipe.Servings != nil && *originalRecipe.Servings > 0 {
				adjustedQuantity.Float64 *= float64(recipe.ServingSize) / float64(*originalRecipe.Servings)
			}

			finalQuantity := adjustedQuantity.Float64
			if finalQuantity <= 0 {
				finalQuantity = 1
			}

			ingredientsWithRecipeIDs = append(ingredientsWithRecipeIDs, IngredientWithRecipeID{
				RecipeID:     recipe.RecipeID,
				IngredientID: ing.IngredientID,
				Name:         ing.Name,
				Quantity:     finalQuantity,
				Unit:         ing.Unit.String,
				Group:        ing.Group.String,
			})
		}
	}

	var markdownTable strings.Builder
	markdownTable.WriteString("| Recipe ID | Ingredient ID | Name | Quantity | Unit | Group |\n")
	markdownTable.WriteString("|-----------|---------------|------|----------|------|-------|\n")

	for _, ingredient := range ingredientsWithRecipeIDs {
		markdownTable.WriteString(fmt.Sprintf("| %s | %d | %s | %.2f | %s | %s |\n",
			ingredient.RecipeID,
			ingredient.IngredientID,
			ingredient.Name,
			ingredient.Quantity,
			ingredient.Unit,
			ingredient.Group,
		))
	}

	aiResponse, err := app.aiClient.GeneratePlanGroceryList(ctx, markdownTable.String())
	if err != nil {
		return plans.GroceryList{}, err
	}

	groceryList := plans.GroceryList{
		GroceryItems: make([]plans.GroceryListItem, len(aiResponse.GroceryListAggregatedIngredients)),
	}

	for i, aiItem := range aiResponse.GroceryListAggregatedIngredients {
		groceryList.GroceryItems[i] = plans.GroceryListItem{
			ID:       aiItem.ID,
			AlikeID:  aiItem.AlikeID,
			RecipeID: aiItem.RecipeID,
			Name:     aiItem.Name,
			Quantity: aiItem.Quantity,
			Unit:     aiItem.Unit,
			Category: aiItem.Category,
			IsMarked: aiItem.IsMarked,
		}
	}

	return groceryList, nil
}

func (a *App) DeleteGroceryListItem(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")
	itemID := c.Param("item_id")

	itemIDInt, err := strconv.Atoi(itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID").SetInternal(err)
	}

	err = plans.DeleteGroceryListItem(ctx, a.db, planID, itemIDInt)
	if err != nil {
		if err == plans.ErrPlanNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "grocery list item not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete grocery list item").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

type groceryListItemUpdateRequest struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Category string  `json:"category"`
}

func (a *App) UpdateGroceryListItem(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")
	itemIDString := c.Param("item_id")

	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid item ID").SetInternal(err)
	}

	var req groceryListItemUpdateRequest
	err = web.ValidateRequest(c, web.ContentTypeApplicationJSON, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = plans.UpdateGroceryListItem(ctx, a.db, planID, itemID, req.Name, req.Quantity, req.Unit, req.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update grocery list item").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) UpdateGroceryListItemMark(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")
	itemIDString := c.Param("item_id")

	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not convert itemID to an int").SetInternal(err)
	}

	err = plans.UpdateGroceryListItemMark(ctx, a.db, planID, itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not toggle grocery list item marked status").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type groceryListItemCreateRequest struct {
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	Unit          string  `json:"unit"`
	Category      string  `json:"category"`
	IsUserCreated bool    `json:"isUserCreated"`
}

func (a *App) CreateGroceryListItem(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	var req groceryListItemCreateRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = plans.CreateGroceryListItem(ctx, a.db, planID, req.Name, req.Quantity, req.Unit, req.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create grocery list item").SetInternal(err)
	}

	return c.NoContent(http.StatusCreated)
}

type updateCategoryOrderRequest struct {
	CategoryOrder []string `json:"categoryOrder"`
}

func (a *App) UpdateCategoryOrder(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	var req updateCategoryOrderRequest
	err := web.ValidateRequest(c, web.ContentTypeApplicationJSON, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = plans.UpdateCategoryOrder(ctx, a.db, planID, req.CategoryOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update category order").SetInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

func (a *App) GetCategoryOrder(c echo.Context) error {
	ctx := c.Request().Context()

	planID := c.Param("id")

	categoryOrder, err := plans.GetCategoryOrderByPlanID(ctx, a.db, planID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get category order").SetInternal(err)
	}

	return c.JSON(http.StatusOK, map[string][]string{"categoryOrder": categoryOrder})
}
