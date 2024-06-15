package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/api"
	"github.com/kitchens-io/kitchens-api/internal/db"
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
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
	}

	// Auth0 Authentication
	Auth0 struct {
		Domain   string        `required:"true" envconfig:"AUTH0_DOMAIN"`
		Audience string        `required:"true" envconfig:"AUTH0_AUDIENCE"`
		CacheTTL time.Duration `default:"5m" envconfig:"AUTH0_CACHE_TTL"`
	}

	// S3 Object Storage configurations
	S3 struct {
		Host string `envconfig:"S3_LOCAL_HOST"`
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

	// AWS Configurations
	//===========================================
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("could not load aws configurations")
	}

	// Handle database migrations and connections.
	// ===========================================
	dsn := db.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	if err := db.Migrate("file://migrations", dsn); err != nil {
		log.Fatal().Err(err).Msg("could migrate database")
	}

	primaryDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start database")
	}
	defer primaryDB.Close()

	// Handle Object storage
	// ==========================
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3.Host != "" {
			o.BaseEndpoint = &cfg.S3.Host
			o.UsePathStyle = true
		}
	})

	// Handle application server.
	// ==========================

	// Create an Auth0 validator.
	validator, err := auth.NewValidator(cfg.Auth0.Domain, cfg.Auth0.Audience, cfg.Auth0.CacheTTL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create jwt validator")
	}

	// Create an API instance.
	app := api.Create(api.CreateInput{
		DB:            primaryDB,
		S3:            s3Client,
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
