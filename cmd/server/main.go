package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/api"
	"github.com/kitchens-io/kitchens-api/pkg/auth"

	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	Debug           bool          `default:"false" envconfig:"DEBUG"`
	Port            string        `default:":1313" envconfig:"PORT"`
	ShutdownTimeout time.Duration `default:"10s" envconfig:"SHUTDOWN_TIMEOUT"`
	Env             string        `default:"dev" envconfig:"APP_ENV"`

	// Database configurations
	// DB struct {
	// 	User string `default:"kitchens_api_svc" envconfig:"DB_USER"`
	// 	Pass string `default:"password" envconfig:"DB_PASS"`
	// 	Host string `default:"localhost:5432" envconfig:"DB_HOST"`
	// 	Name string `default:"kitchens_app" envconfig:"DB_NAME"`
	// 	SSL  bool   `default:"false" envconfig:"DB_SSL"`
	// }

	// Auth0 Authentication
	Auth0 struct {
		Domain   string        `required:"true" envconfig:"AUTH0_DOMAIN"`
		Audience string        `required:"true" envconfig:"AUTH0_AUDIENCE"`
		CacheTTL time.Duration `default:"5m" envconfig:"AUTH0_CACHE_TTL"`
	}
}

func main() {

	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load()

	// Log as JSON instead of the default ASCII formatter.
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Parse environemnt variables into the configuration struct.
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("could not parse application config")
	}

	// Handle database migrations and connections.
	// ===========================================
	// dsn := db.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name, cfg.DB.SSL)
	// if err := db.Migrate("file://migrations", dsn); err != nil {
	// 	log.Fatal(err)
	// }

	// primaryDB, err := db.Connect(dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer primaryDB.Close()

	// Create an Auth0 validator.
	validator, err := auth.NewValidator(cfg.Auth0.Domain, cfg.Auth0.Audience, cfg.Auth0.CacheTTL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create jwt validator")
	}

	// Handle application server.
	// ==========================
	app := api.Create(api.CreateInput{
		AuthValidator: validator,
	})

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := app.API.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("shutting down server")
		}
	}()

	log.Info().
		Interface("config", cfg).
		Msg("kitchens-api started")

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := app.API.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("could not shutdown server")
	}

	log.Info().Msg("kitchens-api gracefully shutdown")
}
