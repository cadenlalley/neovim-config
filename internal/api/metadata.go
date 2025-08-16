package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kitchens-io/kitchens-api/internal/ai"
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
func (a *App) extractRecipeMeta(ctx context.Context, recipeID string) (*ai.RecipeMetaResponseSchema, error) {
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

	result, metrics, err := a.aiClient.ExtractRecipeMetaFromText(ctx, string(recipeJSON))
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("producer", "ai_metadata").
		Str("source", recipeID).
		Str("method", "extractRecipeMeta").
		Interface("metrics", metrics).
		Msg("extracted recipe metadata")

	// TODO: Temporary Backfill recipe metadata, only update if no data has been filled out.
	// This should indicate that the frontend hasn't implemented difficulty, course, or class.
	err = a.backfillRecipeTags(ctx, &recipe, &result)
	if err != nil {
		return &result, err
	}

	// TODO: Temporary Backfill recipe step ingredient associations, only update if no data has been filled out.
	// This should indicate that the frontend hasn't implemented step ingredient associations.
	err = a.backfillRecipeStepIngredients(ctx, &recipe, &result)
	if err != nil {
		return &result, err
	}

	// TODO: Temporary backfill for existing recipes, only if no tags exist.
	recipeHasTags, err := recipes.RecipeHasTags(ctx, a.db, recipeID)
	if err != nil {
		return &result, err
	}

	// If recipe has tags already, do not backfill.
	if recipeHasTags {
		log.Info().Str("recipe_id", recipeID).Msg("recipe already has tags, skipping backfill")
		return &result, nil
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
		return &result, txErr
	}

	return &result, nil
}

// TODO: Temporary backfill functionality
func (a *App) backfillRecipeTags(ctx context.Context, r *recipes.Recipe, result *ai.RecipeMetaResponseSchema) error {
	if r.Difficulty == 0 && !r.Course.Valid && !r.Class.Valid && !r.Cuisine.Valid {
		// uppercase first letter
		result.Cuisine = strings.Title(result.Cuisine)

		err := recipes.BackfillRecipeTags(ctx, a.db, recipes.BackfillRecipeTagsInput{
			RecipeID:   r.RecipeID,
			Difficulty: result.Difficulty,
			Course:     null.StringFrom(result.Course),
			Class:      null.StringFrom(result.Class),
			Cuisine:    null.StringFrom(result.Cuisine),
		})
		if err != nil {
			return err
		}
	} else {
		log.Info().Str("recipe_id", r.RecipeID).Msg("recipe already has metadata, skipping backfill")
	}
	return nil
}

// TODO: Temporary backfill functionality
func (a *App) backfillRecipeStepIngredients(ctx context.Context, r *recipes.Recipe, result *ai.RecipeMetaResponseSchema) error {
	var recipeStepsFound bool
	for _, step := range r.Steps {
		if len(step.IngredientIDs) > 0 {
			recipeStepsFound = true
			break
		}
	}

	if !recipeStepsFound {
		for _, step := range result.StepIngredients {
			if len(step.IngredientIDs) == 0 {
				continue
			}

			err := recipes.BackfillRecipeStepIngredients(ctx, a.db, recipes.BackfillRecipeStepIngredientsInput{
				RecipeID:      r.RecipeID,
				StepID:        step.StepID,
				IngredientIDs: step.IngredientIDs,
			})
			if err != nil {
				return err
			}
		}
	} else {
		log.Info().Str("recipe_id", r.RecipeID).Msg("recipe already has step ingredients, skipping backfill")
	}

	return nil
}
