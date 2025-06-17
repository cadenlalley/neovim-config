package recipes

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v4"
)

var ValidCourses = map[string]interface{}{
	"breakfast": nil,
	"brunch":    nil,
	"lunch":     nil,
	"dinner":    nil,
	"dessert":   nil,
	"supper":    nil,
}

var ValidClasses = map[string]interface{}{
	"main":      nil,
	"side":      nil,
	"snack":     nil,
	"beverage":  nil,
	"dessert":   nil,
	"dip":       nil,
	"soup":      nil,
	"appetizer": nil,
}

var ValidSort = map[string]interface{}{
	"top":  nil,
	"new":  nil,
	"best": nil,
}

type SearchRecipeInput struct {
	Query         string
	KitchenID     null.String
	Course        null.String
	Class         null.String
	Cuisine       null.String
	MaxDifficulty int
	MinRating     int
	MaxTime       int
	OrderBy       null.String
	Limit         uint64
	Offset        uint64
}

func (i *SearchRecipeInput) Validate() error {
	if i.MaxDifficulty > 5 {
		return fmt.Errorf("invalid value for maxDifficulty: '%d'", i.MaxDifficulty)
	}

	if i.MinRating > 5 {
		return fmt.Errorf("invalid value for minRating: '%d'", i.MinRating)
	}

	if i.Course.Valid {
		if _, ok := ValidCourses[i.Course.String]; !ok {
			return fmt.Errorf("invalid value for course: '%s'", i.Course.String)
		}
	}

	if i.Class.Valid {
		if _, ok := ValidClasses[i.Class.String]; !ok {
			return fmt.Errorf("invalid value for class: '%s'", i.Class.String)
		}
	}

	if i.OrderBy.Valid {
		if _, ok := ValidSort[i.OrderBy.String]; !ok {
			return fmt.Errorf("invalid value for orderBy: '%s'", i.OrderBy.String)
		}
	}

	return nil
}

func SearchRecipe(ctx context.Context, store Store, input SearchRecipeInput) ([]SearchResult, error) {
	recipes := make([]SearchResult, 0)

	if err := input.Validate(); err != nil {
		return nil, err
	}

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
		Column("MATCH(r.summary) AGAINST (? IN NATURAL LANGUAGE MODE) AS summary_score", input.Query).
		From("recipes r").
		LeftJoin("(SELECT recipe_id, COUNT(*) as review_count, AVG(rating) as review_rating FROM recipe_reviews GROUP BY recipe_id) rr ON r.recipe_id = rr.recipe_id").
		Where(sq.Expr("MATCH(r.recipe_name, r.summary) AGAINST (? IN NATURAL LANGUAGE MODE)", input.Query)).
		Where(sq.Eq{"r.deleted_at": nil})

	// Apply optional filters
	if input.KitchenID.Valid {
		base = base.Where(sq.Eq{"r.kitchen_id": input.KitchenID.String})
	}

	if input.Course.Valid {
		base = base.Where(sq.Eq{"LOWER(r.course)": strings.ToLower(input.Course.String)})
	}

	if input.Class.Valid {
		base = base.Where(sq.Eq{"LOWER(r.class)": strings.ToLower(input.Class.String)})
	}

	if input.Cuisine.Valid {
		base = base.Where(sq.Eq{"LOWER(r.cuisine)": strings.ToLower(input.Cuisine.String)})
	}

	if input.MaxDifficulty > 0 {
		base = base.Where(sq.LtOrEq{"r.difficulty": input.MaxDifficulty})
	}

	if input.MinRating > 0 {
		base = base.Where(sq.GtOrEq{"rr.review_rating": input.MinRating})
	}

	if input.MaxTime > 0 {
		base = base.Where(sq.LtOrEq{"r.prep_time + r.cook_time": input.MaxTime})
	}

	// Ordering options
	switch input.OrderBy.String {
	case "top":
		base = base.OrderBy("review_rating DESC, review_count DESC")
	case "new":
		base = base.OrderBy("r.created_at DESC")
	case "best":
	default:
		base = base.OrderBy("name_score DESC, summary_score DESC, review_rating DESC, review_count DESC")
	}

	// Apply limit and offset
	base = base.Limit(input.Limit).Offset(input.Offset)

	query, args, err := base.ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}
