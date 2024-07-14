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
	dsn := mysql.DSN(cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	db, err := mysql.Connect(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start database")
	}
	defer db.Close()

	// Reset the migrations database to 0 to force a migration.
	_, err = db.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE TRUE", *cfg.Migrations.Fixtures))
	if err != nil {
		log.Fatal().Err(err).Msg("could not drop migrations database")
	}

	if err := mysql.Migrate("file://fixtures", dsn, cfg.Migrations.Fixtures); err != nil {
		log.Fatal().Err(err).Msg("could not migrate fixtures for database")
	}

	log.Info().Msg("kitchens-api fixtures generated")
}
