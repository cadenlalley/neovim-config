package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/pkg/sdk"
	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	Source string `default:"https://api.kitchens-app.com"`
	JWT    string `required:"true" envconfig:"JWT"`
}

func main() {
	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load()

	// Parse environemnt variables into the configuration struct.
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("could not parse application config")
	}

	// Create kitchens SDK client.
	// =========================================
	kitchensSDK := sdk.NewClient(sdk.NewClientOptions{
		Host: cfg.Source,
	})
	kitchensSDK.SetAuth(cfg.JWT)

	// Fill the recipe ID list for backfilling.
	recipeIDs := []string{}
	if len(recipeIDs) == 0 {
		log.Fatal().Msg("no recipe IDs provided")
		return
	}

	var success int
	var failed int
	for _, recipeID := range recipeIDs {
		log.Info().Msgf("\n====\nprocessing recipe %s\n====\n", recipeID)
		res, err := kitchensSDK.AdminCreateRecipeMetadata(recipeID)
		if err != nil {
			log.Error().Err(err).Msg("could not create recipe metadata")
			failed++
			continue
		}
		log.Info().Interface("tags", res).Msgf("created recipe metadata for recipe %s", recipeID)
		success++
	}

	log.Info().Int("success", success).Int("failed", failed).Msg("completed")
}
