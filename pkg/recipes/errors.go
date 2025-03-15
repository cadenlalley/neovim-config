package recipes

import (
	"errors"
)

var (
	// Recipes
	ErrRecipeNotFound = errors.New("recipe not found")

	// Reviews
	ErrDuplicateReview = errors.New("duplicate review")
	ErrReviewNotFound  = errors.New("review not found")
)
