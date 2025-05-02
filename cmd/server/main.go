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
	"github.com/kitchens-io/kitchens-api/internal/ai"
	"github.com/kitchens-io/kitchens-api/internal/api"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/search"
	"github.com/kitchens-io/kitchens-api/pkg/auth"

	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	Debug           bool          `default:"false" envconfig:"DEBUG"`
	Port            string        `default:"1313" envconfig:"PORT"`
	ShutdownTimeout time.Duration `default:"10s" envconfig:"SHUTDOWN_TIMEOUT"`
	Env             string        `default:"development" envconfig:"APP_ENV"`

	// Database configurations
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
	}

	// Database migrations
	Migrations struct {
		Schema *string `default:"schema_migrations" envconfig:"MIGRATIONS_SCHEMAS"`
	}

	// Auth0 Authentication
	Auth0 struct {
		Domain   string        `required:"true" envconfig:"AUTH0_DOMAIN"`
		Audience string        `required:"true" envconfig:"AUTH0_AUDIENCE"`
		CacheTTL time.Duration `default:"5m" envconfig:"AUTH0_CACHE_TTL"`
	}

	// S3 Object Storage configurations
	S3 struct {
		Host        string `envconfig:"S3_LOCAL_HOST"`
		MediaBucket string `required:"true" envconfig:"S3_MEDIA_BUCKET"`
	}

	CDN struct {
		Host string `required:"true" envconfig:"CDN_HOST"`
	}

	// OpenAI
	OpenAI struct {
		Host  string `required:"true" envconfig:"OPENAI_HOST"`
		Token string `required:"true" envconfig:"OPENAI_TOKEN"`
	}

	// Brave
	Brave struct {
		Host  string `required:"true" envconfig:"BRAVE_HOST"`
		Token string `required:"true" envconfig:"BRAVE_TOKEN"`
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
	dsn := mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	if err := mysql.Migrate("file://migrations", dsn, cfg.Migrations.Schema); err != nil {
		log.Fatal().Err(err).Msg("could not migrate schemas for database")
	}

	db, err := mysql.Connect(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start database")
	}
	defer db.Close()

	// Handle Object storage
	// ==========================
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3.Host != "" {
			o.BaseEndpoint = &cfg.S3.Host
			o.UsePathStyle = true
		}
	})

	fileManager := media.NewS3FileManager(s3Client, cfg.S3.MediaBucket)

	// Handle AI client
	// ==========================
	aiClient := ai.NewClient(cfg.OpenAI.Token, cfg.OpenAI.Host)

	// Handle Search client
	// ==========================
	searchClient := search.NewClient(cfg.Brave.Token, cfg.Brave.Host)

	// Handle application server.
	// ==========================

	// Create an Auth0 validator.
	validator, err := auth.NewValidator(cfg.Auth0.Domain, cfg.Auth0.Audience, cfg.Auth0.CacheTTL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create jwt validator")
	}

	// Create an API instance.
	app := api.Create(api.CreateInput{
		DB:            db,
		FileManager:   fileManager,
		AuthValidator: validator,
		Env:           cfg.Env,
		CDNHost:       cfg.CDN.Host,
		AIClient:      aiClient,
		SearchClient:  searchClient,
	})

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := app.API.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
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
