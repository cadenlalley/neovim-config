package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/pkg/recipes"
	"github.com/kitchens-io/kitchens-api/pkg/tags"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

// extractRecipeMetaBackground wrapper for extractRecipeMeta that should be run as a goroutine.
func (a *App) extractRecipeMetaBackground(recipeID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := a.extractRecipeMeta(ctx, recipeID)
	if err != nil {
		log.Error().
			Str("recipe_id", recipeID).
			Err(err).
			Msg("could not extract recipe metadata")
		return
	}
}

// extractRecipeMeta extracts metadata for a recipe.
func (a *App) extractRecipeMeta(ctx context.Context, recipeID string) ([]tags.Tag, error) {
	recipe, err := recipes.GetRecipeByID(ctx, a.db, recipeID)
	if err != nil {
		if err == recipes.ErrRecipeNotFound {
			return nil, err
		}
		return nil, err
	}

	steps, err := recipes.GetRecipeStepsByRecipeID(ctx, a.db, recipeID)
	if err != nil {
		return nil, err
	}

	notes, err := recipes.GetRecipeNotesByRecipeID(ctx, a.db, recipe.RecipeID)
	if err != nil {
		return nil, err
	}

	for i, step := range steps {
		for _, note := range notes {
			if step.StepID == note.StepID {
				step.Note = note.Note
			}
		}

		steps[i] = step
	}

	recipe.Steps = steps

	ingredients, err := recipes.GetRecipeIngredientsByRecipeID(ctx, a.db, recipe.RecipeID)
	if err != nil {
		return nil, err
	}

	recipe.Ingredients = ingredients

	recipeJSON, err := json.Marshal(recipe)
	if err != nil {
		return nil, err
	}

	result, err := a.aiClient.ExtractRecipeMetaFromText(ctx, string(recipeJSON))
	if err != nil {
		return nil, err
	}

	// TODO: Temporary Backfill recipe metadata, only update if no data has been filled out.
	// This should indicate that the frontend hasn't implemented difficulty, course, or class.
	if recipe.Difficulty == 0 && !recipe.Course.Valid && !recipe.Class.Valid && !recipe.Cuisine.Valid {
		// uppercase first letter
		result.Cuisine = strings.Title(result.Cuisine)

		err = recipes.BackfillRecipeTags(ctx, a.db, recipes.BackfillRecipeTagsInput{
			RecipeID:   recipeID,
			Difficulty: result.Difficulty,
			Course:     null.StringFrom(result.Course),
			Class:      null.StringFrom(result.Class),
			Cuisine:    null.StringFrom(result.Cuisine),
		})
		if err != nil {
			return nil, err
		}
	}

	// Convert to tags
	var tagResponse []tags.Tag
	for _, tag := range result.Tags {
		tagResponse = append(tagResponse, tags.Tag{
			Type:  tag.Type,
			Value: strings.Replace(tag.Value, " ", "-", -1),
		})
	}

	var createdTags []tags.Tag
	txErr := mysql.Transaction(ctx, a.db, func(tx *sqlx.Tx) error {
		createdTags, err = tags.CreateTags(ctx, tx, tagResponse)
		if err != nil {
			return err
		}

		tagIDs := make([]int, len(createdTags))
		for i, tag := range createdTags {
			tagIDs[i] = tag.TagID
		}

		err = recipes.CreateRecipeTags(ctx, tx, recipeID, tagIDs)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return createdTags, nil
}
