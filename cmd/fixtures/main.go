package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
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

	log.Info().Msg("kitchens-api fixtures generated")
}
