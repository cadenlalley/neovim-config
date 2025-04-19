package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DSN converts a database configuration to a connection string.
func DSN(user, pass, host, name string) string {
	return fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, name)
}

// Connect returns a pointer to an initialized database.
func Connect(dsn string) (*sqlx.DB, error) {
	parts := strings.Split(dsn, "mysql://")

	d, err := sqlx.Connect("mysql", parts[1])
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Migrate applies database migrations from the migrationsPath to the
// database specified in the dsn string.
func Migrate(migrationsPath, dsn string, migrationsTable *string) error {
	table := "schema_migrations"
	if migrationsTable != nil {
		table = *migrationsTable
	}

	input := fmt.Sprintf("%s&multiStatements=true&x-migrations-table=%s", dsn, table)
	m, err := migrate.New(migrationsPath, input)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// ResetFixtures will reset the database by dropping the migrations table
// and re-applying the migrations from the migrationsPath.
func ResetFixtures(migrationsPath, dsn string, migrationsTable *string) error {
	table := "fixtures_migrations"
	if migrationsTable != nil {
		table = *migrationsTable
	}

	input := fmt.Sprintf("%s&multiStatements=true&x-migrations-table=%s", dsn, table)
	m, err := migrate.New(migrationsPath, input)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func Transaction(ctx context.Context, conn *sqlx.DB, f func(tx *sqlx.Tx) error) error {
	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := f(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
