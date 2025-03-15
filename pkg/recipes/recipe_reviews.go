package recipes

import (
	"context"
	"database/sql"
	"strings"

	"gopkg.in/guregu/null.v4"
)

type CreateReviewInput struct {
	ReviewID          string
	RecipeID          string
	ReviewerKitchenID string
	Description       null.String
	Rating            float64
	Media             null.String
}

func CreateReview(ctx context.Context, store Store, input CreateReviewInput) (Review, error) {
	// Handle nullable values.
	if input.Description.String == "" {
		input.Description = null.NewString(input.Description.String, false)
	}
	if input.Media.String == "" {
		input.Media = null.NewString(input.Media.String, false)
	}

	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_reviews (
			review_id,
			recipe_id,
			reviewer_kitchen_id,
			review_description,
			rating,
			media_path
		) VALUES (?, ?, ?, ?, ?, ?)
	`, input.ReviewID, input.RecipeID, input.ReviewerKitchenID, input.Description, input.Rating, input.Media)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return Review{}, ErrDuplicateReview
		}
		return Review{}, err
	}

	return GetReviewByID(ctx, store, input.ReviewID)
}

type UpdateReviewInput struct {
	ReviewID    string
	Description null.String
	Rating      float64
	Media       null.String
}

func UpdateReview(ctx context.Context, store Store, input UpdateReviewInput) (Review, error) {
	// Handle nullable values.
	if input.Description.String == "" {
		input.Description = null.NewString(input.Description.String, false)
	}
	if input.Media.String == "" {
		input.Media = null.NewString(input.Media.String, false)
	}

	_, err := store.ExecContext(ctx, `
		UPDATE recipe_reviews SET
			review_description = ?,
			rating = ?,
			media_path = ?
		WHERE review_id = ?;
	`, input.Description, input.Rating, input.Media, input.ReviewID)
	if err != nil {
		return Review{}, err
	}

	return GetReviewByID(ctx, store, input.ReviewID)
}

func DeleteReview(ctx context.Context, store Store, reviewID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipe_reviews WHERE review_id = ?;
	`, reviewID)

	if err != nil {
		return err
	}

	return nil
}

func GetReviewByID(ctx context.Context, store Store, reviewID string) (Review, error) {
	var review Review
	err := store.QueryRowxContext(ctx, `
		SELECT * FROM recipe_reviews WHERE review_id = ?;
	`, reviewID).StructScan(&review)
	if err != nil {
		if err == sql.ErrNoRows {
			return Review{}, ErrReviewNotFound
		}
		return Review{}, err
	}

	return review, nil
}

func ListReviewsByRecipeID(ctx context.Context, store Store, recipeID string, activeKitchenID string) ([]Review, error) {
	// Retrieves recipe reviews with enriched metadata, joining kitchen and account details.
	// Calculates total review likes and checks if the current kitchen has liked each review.
	rows, err := store.QueryxContext(ctx, `
		WITH reviews_enriched AS (
			SELECT
				rr.*,
						concat(a.first_name, ' ', a.last_name) as reviewer_name,
				k.avatar as reviewer_avatar
			FROM recipe_reviews rr
				LEFT JOIN kitchens k ON rr.reviewer_kitchen_id = k.kitchen_id
				LEFT JOIN accounts a ON k.account_id = a.account_id
			WHERE rr.recipe_id = ?
		)
		SELECT
			re.*,
			COALESCE(rrl.total_likes, 0) as total_likes,
			COALESCE(rrl.liked, 0) as liked
		FROM reviews_enriched re
			LEFT JOIN (
				SELECT
					review_id,
					count(*) as total_likes,
					SUM(CASE WHEN kitchen_id = ? THEN 1 ELSE 0 END) as liked
				FROM recipe_review_likes GROUP BY review_id
			) AS rrl ON re.review_id = rrl.review_id
			WHERE recipe_id = ?;
	`, recipeID, activeKitchenID, recipeID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviews []Review

	for rows.Next() {
		var review Review
		if err := rows.StructScan(&review); err != nil {
			return reviews, err
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return reviews, err
	}

	return reviews, nil
}

func GetReviewSummaryByRecipeID(ctx context.Context, store Store, recipeID string) (ReviewSummary, error) {
	var reviewSummary ReviewSummary
	err := store.QueryRowxContext(ctx, `
		SELECT
			count(*) as total,
			COALESCE(avg(rating), 0) as average,
			COALESCE(SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END), 0) as rating_1,
			COALESCE(SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END), 0) as rating_2,
			COALESCE(SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END), 0) as rating_3,
			COALESCE(SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END), 0) as rating_4,
			COALESCE(SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END), 0) as rating_5
		FROM recipe_reviews WHERE recipe_id = ?;
	`, recipeID).StructScan(&reviewSummary)
	if err != nil {
		if err == sql.ErrNoRows {
			return ReviewSummary{}, ErrReviewNotFound
		}
		return ReviewSummary{}, err
	}

	reviewSummary.Ratings = map[int]int{
		1: reviewSummary.Rating_1,
		2: reviewSummary.Rating_2,
		3: reviewSummary.Rating_3,
		4: reviewSummary.Rating_4,
		5: reviewSummary.Rating_5,
	}

	return reviewSummary, nil
}

func LikeReview(ctx context.Context, store Store, reviewID string, kitchenID string) error {
	_, err := store.ExecContext(ctx, `
		INSERT INTO recipe_review_likes (review_id, kitchen_id)
		VALUES (?, ?);
	`, reviewID, kitchenID)
	if err != nil {
		return err
	}

	return nil
}

func UnlikeReview(ctx context.Context, store Store, reviewID string, kitchenID string) error {
	_, err := store.ExecContext(ctx, `
		DELETE FROM recipe_review_likes WHERE review_id = ? AND kitchen_id = ?;
	`, reviewID, kitchenID)
	if err != nil {
		return err
	}

	return nil
}
