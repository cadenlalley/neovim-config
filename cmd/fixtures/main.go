package main

import (
	"context"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kitchens-io/kitchens-api/internal/media"
	"github.com/kitchens-io/kitchens-api/internal/mysql"

	"github.com/rs/zerolog/log"
)

// Application Configuration
type AppConfig struct {
	// Database configurations
	DB struct {
		User string `required:"true" envconfig:"DB_USER"`
		Pass string `required:"true" envconfig:"DB_PASS"`
		Host string `required:"true" envconfig:"DB_HOST"`
		Name string `required:"true" envconfig:"DB_NAME"`
	}
	Migrations struct {
		Fixtures *string `default:"fixtures_migrations" envconfig:"MIGRATIONS_FIXTURES"`
		Schema   *string `default:"schema_migrations" envconfig:"MIGRATIONS_SCHEMAS"`
	}
	S3 struct {
		Host        string `envconfig:"S3_LOCAL_HOST"`
		MediaBucket string `required:"true" envconfig:"S3_MEDIA_BUCKET"`
	}
}

func main() {

	// .env file is optional, so ignore the error returned from loading.
	_ = godotenv.Load()

	// Parse environemnt variables into the configuration struct.
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal().Err(err).Msg("could not parse application config")
	}

	// Handle database migrations and connections.
	// ===========================================

	// 1. Create an initial connection just to th DB with no database selected.
	//    This will allow for the dropping of the DB.
	db1, err := mysql.Connect(mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, ""))
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}
	defer db1.Close()

	_, err = db1.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DB.Name))
	if err != nil {
		log.Fatal().Err(err).Msg("could not drop database")
	}

	_, err = db1.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DB.Name))
	if err != nil {
		log.Fatal().Err(err).Msg("could not create database")
	}

	// 2. Create a connection directly to the DB with a database selected.
	//    Run the migrations for the schemas and then apply fixtures.
	dsn := mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)
	db, err := mysql.Connect(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start database")
	}
	defer db.Close()

	if err := mysql.Migrate("file://migrations", dsn, cfg.Migrations.Schema); err != nil {
		log.Fatal().Err(err).Msg("could not migrate schemas for database")
	}

	if err := mysql.Migrate("file://fixtures", dsn, cfg.Migrations.Fixtures); err != nil {
		log.Fatal().Err(err).Msg("could not migrate fixtures for database")
	}

	// Handle file uploads for exisiting resources.
	// ===========================================
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal().Err(err).Msg("could not load aws configurations")
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3.Host != "" {
			o.BaseEndpoint = &cfg.S3.Host
			o.UsePathStyle = true
		}
	})

	fileManager := media.NewS3FileManager(s3Client, cfg.S3.MediaBucket)

	// 1. Iterate over the fixtures/media directory and upload all files.
	media := map[string]string{
		// Kitchens
		"kitchen_avatar_sammycooks.png": "kitchens/ktc_2jEx1e1esA5292rBisRGuJwXc14/2pR9Azi4FC8IkJO3J6ZvjUqgz5Z.png",
		"kitchen_cover_sammycooks.png":  "kitchens/ktc_2jEx1e1esA5292rBisRGuJwXc14/2pR9B1iDNJ7EXddtd1mDzCvXaJt.png",
		"kitchen_avatar_bbq_bill.jpg":   "kitchens/ktc_2jEx1j3CVPIIAaOwGIORKqHfK89/2qwGAZKuTbHO96b6VXZohC6VBes.png",

		// Recipes
		"recipe_pumpkin_pie.png": "recipes/rcp_2jbgfAMKOCnKrWQroRBkXPIRI6T/2pR9B2cIFxj82GDTDB44lpMzYHu.png",

		// Folders
		"folder_breakfast.png":     "folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6BliLuDCHYJKq7Eqkb9l55bS.png",
		"folder_healthy_lunch.png": "folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6Bouca6LsCl1fXHDK9p264bi.png",
		"folder_mediterranean.png": "folders/fld_2pPgQjn08dQzr5vjSk8WYSBTATo/2pR6BrWvEtBW5uBwoNc213qALBc.png",
	}

	for sourceFile, targetKey := range media {
		path := fmt.Sprintf("fixtures/media/%s", sourceFile)
		file, err := os.Open(path)
		if err != nil {
			log.Fatal().Err(err).Msgf("could not open file: %s", sourceFile)
		}
		defer file.Close()

		key := fmt.Sprintf("uploads/%s", targetKey)
		if err := fileManager.Upload(context.TODO(), file, key); err != nil {
			log.Fatal().Err(err).Msgf("could not upload file: %s", sourceFile)
		}
	}

	log.Info().Msg("kitchens-api fixtures generated")
}
