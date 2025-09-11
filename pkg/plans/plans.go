package plans

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}

type CreatePlanInput struct {
	ID             string
	OwnerAccountID string
	StartDate      time.Time
	EndDate        time.Time
}

func CreatePlan(ctx context.Context, store Store, input CreatePlanInput) (Plan, error) {
	findQuery := `
		SELECT meal_plan_id, account_id, start_date, end_date, created_at, updated_at
		FROM meal_plans
		WHERE meal_plan_id = ?;
	`
	createQuery := `
		INSERT INTO meal_plans (meal_plan_id, account_id, start_date, end_date)
		VALUES (?, ?, ?, ?);
	`

	row := store.QueryRowxContext(ctx, findQuery, input.ID)
	if row.Err() != nil {
		return Plan{}, row.Err()
	}

	var plan Plan
	err := row.Scan(&plan.ID, &plan.OwnerAccountID, &plan.StartDate, &plan.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := store.ExecContext(ctx, createQuery, input.ID, input.OwnerAccountID, input.StartDate, input.EndDate)
			if err != nil {
				return Plan{}, err
			}

			plan, err := GetPlanByID(ctx, store, input.ID)
			if err != nil {
				return Plan{}, err
			}

			return plan, nil
		}
		return Plan{}, err
	}

	return plan, nil
}

func GetPlanByID(ctx context.Context, store Store, id string) (Plan, error) {
	query := `
		SELECT meal_plan_id, account_id, start_date, end_date, created_at, updated_at
		FROM meal_plans
		WHERE meal_plan_id = ?;
	`
	row := store.QueryRowxContext(ctx, query, id)

	var plan Plan
	err := row.Scan(&plan.ID, &plan.OwnerAccountID, &plan.StartDate, &plan.EndDate, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Plan{}, ErrPlanNotFound
		}
		return Plan{}, fmt.Errorf("failed to get plan by ID: %w", err)
	}

	return plan, nil
}

func ListPlansByUserID(ctx context.Context, store Store, userID string) ([]Plan, error) {
	query := `
		SELECT meal_plan_id, account_id, start_date, end_date, created_at, updated_at
		FROM meal_plans
		WHERE account_id = ?;
	`
	rows, err := store.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var plan Plan
		err := rows.Scan(&plan.ID, &plan.OwnerAccountID, &plan.StartDate, &plan.EndDate)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func GetPlanByAccountIDAndDateRange(ctx context.Context, store Store, accountID string, startDate time.Time, endDate time.Time) (Plan, error) {
	query := `
		SELECT meal_plan_id, account_id, start_date, end_date, created_at, updated_at
		FROM meal_plans
		WHERE account_id = ? AND start_date >= ? AND end_date <= ?;
	`
	row := store.QueryRowxContext(ctx, query, accountID, startDate, endDate)
	if row.Err() != nil {
		return Plan{}, row.Err()
	}

	var plan Plan
	err := row.Scan(&plan.ID, &plan.OwnerAccountID, &plan.StartDate, &plan.EndDate, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Plan{}, ErrPlanNotFound
		}
		return Plan{}, err
	}

	return plan, nil
}

func GetRecipesByPlanID(ctx context.Context, store Store, planID string) ([]PlanRecipe, error) {
	query := `
		SELECT recipe_id, day_number, serving_size 
		FROM meal_plan_recipes
		WHERE meal_plan_id = ?;
	`
	rows, err := store.QueryxContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var recipes []PlanRecipe
	for rows.Next() {
		var recipe PlanRecipe
		err := rows.Scan(&recipe.RecipeID, &recipe.Day, &recipe.ServingSize)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func GetFullRecipesByPlanID(ctx context.Context, store Store, planID string) ([]FullPlanRecipe, error) {
	query := `
		SELECT 
			mpr.meal_plan_recipe_id,
			r.recipe_id, 
			mpr.day_number, 
			mp.start_date,
			mp.end_date,
			mpr.serving_size, 
			r.kitchen_id, 
			r.recipe_name, 
			r.summary,
	       	r.prep_time, 
			r.cook_time, 
			r.servings, 
			r.difficulty, 
			r.course, 
			r.class, 
			r.cuisine, 
			r.cover,
		    r.source, 
			r.created_at, 
			r.updated_at, 
			r.deleted_at
		FROM meal_plan_recipes mpr
		JOIN recipes r ON mpr.recipe_id = r.recipe_id
		JOIN meal_plans mp ON mpr.meal_plan_id = mp.meal_plan_id
		WHERE mpr.meal_plan_id = ?;
	`
	rows, err := store.QueryxContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fullRecipes []FullPlanRecipe
	for rows.Next() {
		var recipe FullPlanRecipe
		err := rows.StructScan(&recipe)
		if err != nil {
			return nil, err
		}
		fullRecipes = append(fullRecipes, recipe)
	}

	return fullRecipes, nil
}

func AddRecipesToPlan(ctx context.Context, store Store, planID string, recipeID string, dayNumber int, servingSize int) error {
	query := `
		INSERT INTO meal_plan_recipes (meal_plan_id, recipe_id, day_number, serving_size)
		VALUES (?, ?, ?, ?);
	`
	_, err := store.ExecContext(ctx, query, planID, recipeID, dayNumber, servingSize)
	if err != nil {
		return err
	}

	return nil
}

func RemoveRecipeFromPlan(ctx context.Context, store Store, mealPlanRecipeID string) error {
	query := `
		DELETE FROM meal_plan_recipes WHERE meal_plan_recipe_id = ?
	`

	_, err := store.ExecContext(ctx, query, mealPlanRecipeID)
	if err != nil {
		return err
	}

	return nil
}

func GetGroceryListByMealPlanID(ctx context.Context, store Store, planID string) (*GroceryList, error) {
	query := `SELECT item_id, alike_id, recipe_id, name, quantity, unit, category, marked FROM meal_plan_grocery_list_items WHERE meal_plan_id = ?`

	rows, err := store.QueryxContext(ctx, query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groceryListItems := &GroceryList{
		GroceryItems: []GroceryListItem{},
	}

	hasRows := false
	for rows.Next() {
		hasRows = true
		var item GroceryListItem

		err := rows.Scan(
			&item.ID,
			&item.AlikeID,
			&item.RecipeID,
			&item.Name,
			&item.Quantity,
			&item.Unit,
			&item.Category,
			&item.IsMarked,
		)
		if err != nil {
			return nil, err
		}

		groceryListItems.GroceryItems = append(groceryListItems.GroceryItems, item)
	}

	// If no rows found, return nil to trigger AI generation
	if !hasRows {
		return nil, nil
	}

	return groceryListItems, nil
}

func CreateGroceryList(ctx context.Context, store Store, planID string, groceryList GroceryList) error {
	if len(groceryList.GroceryItems) == 0 {
		return nil
	}

	builder := sq.
		StatementBuilder.
		PlaceholderFormat(sq.Question).
		Insert("meal_plan_grocery_list_items").
		Columns("meal_plan_id", "alike_id", "recipe_id", "name", "quantity", "unit", "category")

	for _, groceryItem := range groceryList.GroceryItems {
		builder = builder.Values(planID, groceryItem.AlikeID, groceryItem.RecipeID, groceryItem.Name, groceryItem.Quantity, groceryItem.Unit, groceryItem.Category)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = store.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteGroceryList(ctx context.Context, store Store, planID string) error {
	query := `DELETE FROM meal_plan_grocery_list_items WHERE meal_plan_id = ? AND is_user_created = false`

	_, err := store.ExecContext(ctx, query, planID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteGroceryListItem(ctx context.Context, store Store, planID string, itemID int) error {
	query := `DELETE FROM meal_plan_grocery_list_items WHERE meal_plan_id = ? AND item_id = ?`

	result, err := store.ExecContext(ctx, query, planID, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPlanNotFound
	}

	return nil
}

func UpdateGroceryListItemMark(ctx context.Context, store Store, planID string, itemID int) error {
	query := `UPDATE meal_plan_grocery_list_items SET marked = NOT marked WHERE meal_plan_id = ? AND item_id = ?`

	_, err := store.ExecContext(ctx, query, planID, itemID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateGroceryListItem(ctx context.Context, store Store, planID string, itemID int, name string, quantity float64, unit string, category string) error {
	query := `UPDATE meal_plan_grocery_list_items SET name = ?, quantity = ?, unit = ?, category = ? WHERE meal_plan_id = ? AND item_id = ?`

	_, err := store.ExecContext(ctx, query, name, quantity, unit, category, planID, itemID)
	if err != nil {
		return err
	}

	return nil
}

func CreateGroceryListItem(ctx context.Context, store Store, planID string, name string, quantity float64, unit string, category string) error {
	query := `INSERT INTO meal_plan_grocery_list_items (meal_plan_id, recipe_id, name, quantity, unit, category, is_user_created, marked, alike_id) VALUES (?, '', ?, ?, ?, ?, true, false, 0)`

	_, err := store.ExecContext(ctx, query, planID, name, quantity, unit, category)
	if err != nil {
		return err
	}

	return nil
}

func GetPlanGroceryListIsDirty(ctx context.Context, store Store, planID string) (bool, error) {
	query := `SELECT grocery_list_is_dirty FROM meal_plans WHERE meal_plan_id = ?`

	row := store.QueryRowxContext(ctx, query, planID)

	var isDirty bool
	err := row.Scan(&isDirty)
	if err != nil {
		return true, err
	}

	return isDirty, nil
}

func SetPlanGroceryListIsDirty(ctx context.Context, store Store, planID string, isDirty bool) error {
	query := `UPDATE meal_plans SET grocery_list_is_dirty = ? WHERE meal_plan_id = ?`

	_, err := store.ExecContext(ctx, query, isDirty, planID)
	if err != nil {
		return err
	}

	return nil
}
