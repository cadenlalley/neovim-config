package plans

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
