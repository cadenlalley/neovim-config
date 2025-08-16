package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	comp "github.com/adrg/strutil/metrics"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/pkg/sdk"
	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	JWT    string `default:"test" required:"true" envconfig:"JWT"`
	Source string `default:"http://localhost:1313"`
}

type Evaluation struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	ExpectedFile string `json:"expectedFile"`
}

type EvaluationResult struct {
	Metrics  ai.ResponseMetrics `json:"metrics"`
	Distance float64            `json:"distance"`
}

type ModelMetadata struct {
	// Per Million Token Cost
	InputCost  float64 `json:"inputCost"`
	OutputCost float64 `json:"outputCost"`
}

var Models map[string]ModelMetadata = map[string]ModelMetadata{
	"gpt-4o-mini-2024-07-18": {
		InputCost:  0.15,
		OutputCost: 0.60,
	},
	"gpt-4.1-mini-2025-04-14": {
		InputCost:  0.40,
		OutputCost: 1.60,
	},
	"gpt-4o-2024-11-20": {
		InputCost:  2.50,
		OutputCost: 10.00,
	},
	"gpt-5-nano-2025-08-07": {
		InputCost:  0.05,
		OutputCost: 0.40,
	},
	"gpt-5-mini-2025-08-07": {
		InputCost:  0.25,
		OutputCost: 2.00,
	},
}

func main() {
	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load()

	// Parse environment variables into the configuration struct.
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

	// Load evaluations fromfile
	// =========================================
	var evaluations []Evaluation
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config file")
	}
	json.Unmarshal(b, &evaluations)

	// Run evaluations.
	// =========================================
	results := make(map[string]EvaluationResult)

	for _, evaluation := range evaluations {
		log.Info().Str("name", evaluation.Name).Msg("evaluating")

		recipe, metrics, err := kitchensSDK.ImportRecipeFromURL(evaluation.URL)
		if err != nil {
			log.Err(err).Str("name", evaluation.Name).Msg("could not import recipe from URL")
			continue
		}

		actual, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			log.Err(err).Str("name", evaluation.Name).Msg("could not marshal recipe")
			continue
		}

		dirName := "./actual_" + strings.ReplaceAll(metrics.Model, " ", "_")

		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			err = os.Mkdir(dirName, 0755)
			if err != nil {
				log.Err(err).Str("name", evaluation.Name).Msg("could not create actual directory")
				continue
			}
		}

		err = os.WriteFile(strings.Replace(evaluation.ExpectedFile, "./expected", dirName, 1), actual, 0644)
		if err != nil {
			log.Err(err).Str("name", evaluation.Name).Msg("could not write actual file")
		}

		expected, err := os.ReadFile(evaluation.ExpectedFile)
		if err != nil {
			log.Err(err).Str("name", evaluation.Name).Msg("could not read expected file")
			continue
		}

		comparator := comp.NewLevenshtein()
		comparator.CaseSensitive = false

		distance := comparator.Compare(string(expected), string(actual))
		results[evaluation.Name] = EvaluationResult{
			Metrics:  *metrics,
			Distance: distance,
		}
	}

	fmt.Println("| Source | Model | Prompt Tokens | Completion Tokens | Latency (ms) | Accuracy | Total Cost |")
	fmt.Println("| --- | --- | --- | --- | --- | --- | --- |")
	for k, v := range results {
		promptTokenCost := (float64(v.Metrics.PromptTokens) / 1e6) * Models[v.Metrics.Model].InputCost
		completionTokenCost := (float64(v.Metrics.CompletionTokens) / 1e6) * Models[v.Metrics.Model].OutputCost
		totalCost := promptTokenCost + completionTokenCost

		fmt.Printf("| %s | %s | %d | %d | %d | %f | %f |\n", k, v.Metrics.Model, v.Metrics.PromptTokens, v.Metrics.CompletionTokens, v.Metrics.Latency, v.Distance, totalCost)
	}
}
