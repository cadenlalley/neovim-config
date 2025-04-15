package api

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"
	"github.com/kitchens-io/kitchens-api/internal/openai"
	"github.com/rs/zerolog/log"
)

var testApp *App

// Application Configuration
type AppConfig struct {
	Debug           bool          `default:"false" envconfig:"DEBUG"`
	Port            string        `default:"1313" envconfig:"PORT"`
	ShutdownTimeout time.Duration `default:"10s" envconfig:"SHUTDOWN_TIMEOUT"`
	Env             string        `default:"test" envconfig:"APP_ENV"`

	// Database configurations
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
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
}

func TestMain(m *testing.M) {
	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load(filepath.Join("../../", ".env"))

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

	aiClient := openai.NewOpenAIClient(cfg.OpenAI.Host, cfg.OpenAI.Token, cfg.Debug)

	// Handle application server.
	// ==========================

	// Create an API instance.
	testApp = Create(CreateInput{
		DB:            db,
		FileManager:   fileManager,
		AuthValidator: nil,
		Env:           cfg.Env,
		CDNHost:       cfg.CDN.Host,
		AIClient:      aiClient,
	})

	// Run all tests
	exitCode := m.Run()

	// Exit with the test result code
	os.Exit(exitCode)
}
