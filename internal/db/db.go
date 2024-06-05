package db

import (
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
func Migrate(migrationsPath, dsn string) error {
	m, err := migrate.New(migrationsPath, dsn+"&multiStatements=true")
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// ApplyFixtures will receive a DB and table to reset for tests.
// This function can only be used locally or in tests.
// TODO: Figure out pathing issue for tests.
// func ApplyFixtures(db *sqlx.DB, path, table string) error {
// 	if env := os.Getenv("APP_ENV"); env != "dev" && env != "test" {
// 		return fmt.Errorf("Unable to apply fixtures in environment: %s", env)
// 	}

// 	var fileName string
// 	switch table {
// 	case "users":
// 		fileName = "01_create_users.up.sql"
// 	default:
// 		return fmt.Errorf("Invalid table name supplied '%s'", table)
// 	}

// 	dsql := fmt.Sprintf(`DELETE FROM %s WHERE 1=1;`, table)
// 	_, err := db.Exec(dsql)
// 	if err != nil {
// 		return err
// 	}

// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}

// 	fixture := filepath.Join(pwd, path, "fixtures", fileName)
// 	sql, err := ioutil.ReadFile(fixture)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := db.Exec(string(sql)); err != nil {
// 		return err
// 	}

// 	return nil
// }
