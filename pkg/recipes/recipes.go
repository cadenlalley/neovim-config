package recipes

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type Store interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}

type CreateRecipeInput struct {
	RecipeID   string
	KitchenID  string
	Name       string
	Summary    null.String
	PrepTime   int
	CookTime   int
	Servings   int
	Difficulty int
	Course     null.String
	Class      null.String
	Cuisine    null.String
	Cover      null.String
	Source     null.String
}

func CreateRecipe(ctx context.Context, store Store, input CreateRecipeInput) (Recipe, error) {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipes (
			recipe_id,
			kitchen_id,
			recipe_name,
			summary,
			prep_time,
			cook_time,
			servings,
			difficulty,
			course,
			class,
			cuisine,
			cover,
			source,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);
	`, input.RecipeID, input.KitchenID, input.Name, input.Summary, input.PrepTime, input.CookTime, input.Servings, input.Difficulty, input.Course, input.Class, input.Cuisine, input.Cover, input.Source)

	if err != nil {
		return Recipe{}, err
	}

	return GetRecipeByID(ctx, store, input.RecipeID)
}

type UpdateRecipeInput struct {
	RecipeID   string
	Name       string
	Summary    null.String
	PrepTime   int
	CookTime   int
	Servings   int
	Difficulty int
	Course     null.String
	Class      null.String
	Cuisine    null.String
	Cover      null.String
	Source     null.String
}

func UpdateRecipe(ctx context.Context, store Store, input UpdateRecipeInput) (Recipe, error) {
	_, err := store.ExecContext(ctx, `
		UPDATE recipes
		SET
			recipe_name = ?,
			summary = ?,
			prep_time = ?,
			cook_time = ?,
			servings = ?,
			difficulty = ?,
			course = ?,
			class = ?,
			cuisine = ?,
			cover = ?,
			source = ?
		WHERE
			recipe_id = ?;
	`, input.Name, input.Summary, input.PrepTime, input.CookTime, input.Servings, input.Difficulty, input.Course, input.Class, input.Cuisine, input.Cover, input.Source, input.RecipeID)
	if err != nil {
		return Recipe{}, err
	}

	return GetRecipeByID(ctx, store, input.RecipeID)
}

func GetRecipeByID(ctx context.Context, store Store, recipeID string) (Recipe, error) {
	var recipe Recipe
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM recipes WHERE recipe_id = ?;
	`, recipeID).StructScan(&recipe)

	if err != nil {
		if err == sql.ErrNoRows {
			return Recipe{}, ErrRecipeNotFound
		}
		return Recipe{}, err
	}

	// TODO: Revisit if this is the best way to handle this.
	err = recipe.ComputeValues()
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func ListRecipesByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Recipe, error) {
	recipes := make([]Recipe, 0)

	rows, err := store.QueryxContext(ctx, `
		SELECT
			r.*,
			CASE WHEN rr.review_count IS NOT NULL THEN rr.review_count ELSE 0 END as review_count,
			CASE WHEN rr.review_rating IS NOT NULL THEN rr.review_rating ELSE 0 END as review_rating
		FROM recipes r
			LEFT JOIN (SELECT recipe_id, count(*) as review_count, avg(rating) as review_rating
									FROM recipe_reviews
									GROUP BY recipe_id
								) AS rr ON r.recipe_id = rr.recipe_id
		WHERE (r.kitchen_id = ? OR r.recipe_id IN (SELECT recipe_id FROM recipes_saved WHERE kitchen_id = ?))
		  AND r.deleted_at IS NULL
		ORDER BY r.created_at;
	`, kitchenID, kitchenID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var recipe Recipe
		if err := rows.StructScan(&recipe); err != nil {
			return recipes, err
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return recipes, err
	}

	return recipes, nil
}

func DeleteRecipeByID(ctx context.Context, store Store, recipeID string) error {
	_, err := store.ExecContext(ctx, `
		UPDATE recipes SET deleted_at = CURRENT_TIMESTAMP WHERE recipe_id = ?
	`, recipeID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRecipesByKitchenID(ctx context.Context, store Store, kitchenID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipes WHERE kitchen_id = ?
	`, kitchenID)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Temporary for backfilling recipes with missing difficulty, course, class until frontend sends consistently.
type BackfillRecipeTagsInput struct {
	RecipeID   string
	Difficulty int
	Course     null.String
	Class      null.String
	Cuisine    null.String
}

func BackfillRecipeTags(ctx context.Context, store Store, input BackfillRecipeTagsInput) error {
	_, err := store.ExecContext(ctx, `
		UPDATE recipes
			SET difficulty = ?,
			course = ?,
			class = ?,
			cuisine = ?
		WHERE recipe_id = ?;
	`, input.Difficulty, input.Course, input.Class, input.Cuisine, input.RecipeID)
	if err != nil {
		return err
	}
	return nil
}

type SearchRecipeInput struct {
	Query     string
	KitchenID string
}

func SearchRecipe(ctx context.Context, store Store, input SearchRecipeInput) ([]SearchResult, error) {
	recipes := make([]SearchResult, 0)

	query, args, err := prepareSearchRecipeQuery(input)
	if err != nil {
		return nil, err
	}

	rows, err := store.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var recipe SearchResult
		if err := rows.StructScan(&recipe); err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

// Prepare the search query for recipes.
func prepareSearchRecipeQuery(input SearchRecipeInput) (string, []interface{}, error) {
	base := sq.Select(
		"r.recipe_id",
		"r.recipe_name",
		"r.kitchen_id",
		"r.source",
		"r.cover",
		"COALESCE(rr.review_count, 0) as review_count",
		"COALESCE(rr.review_rating, 0) as review_rating",
	).
		Column("MATCH(r.recipe_name) AGAINST (? IN NATURAL LANGUAGE MODE) * 2.0 AS name_score", input.Query).
		Column("MATCH(r.summary) AGAINST (? IN NATURAL LANGUAGE MODE) * 1.0 AS summary_score", input.Query).
		From("recipes r").
		LeftJoin("(SELECT recipe_id, COUNT(*) as review_count, AVG(rating) as review_rating FROM recipe_reviews GROUP BY recipe_id) rr ON r.recipe_id = rr.recipe_id").
		Where(sq.Expr("MATCH(r.recipe_name, r.summary) AGAINST (? IN NATURAL LANGUAGE MODE)", input.Query)).
		Where(sq.Eq{"r.deleted_at": nil})

	// Apply optional filters
	if input.KitchenID != "" {
		base = base.Where(sq.Eq{"r.kitchen_id": input.KitchenID})
	}

	// Order by name score and summary score
	query := base.OrderBy("name_score DESC, summary_score DESC").Limit(20)
	output, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return output, args, nil
}
