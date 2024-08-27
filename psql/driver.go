package psql

import (
	"fmt"
	"github.com/gflydev/db"
	"github.com/jiveio/fluentsql"
	"github.com/jmoiron/sqlx"
	"os"

	// Autoload driver for PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ========================================================================================
//                                     PostgreSQL Driver
// ========================================================================================

// New initial PostgreSQL driver and register to database manager
func New() *PostgreSQL {
	// Set DBType
	fluentsql.SetDBType(fluentsql.PostgreSQL)

	// Create driver
	return &PostgreSQL{}
}

// PostgreSQL a implement of interface IDatabase for PostgreSQL
type PostgreSQL struct{}

// Load perform DB connection to PostgreSQL database.
func (d *PostgreSQL) Load() (*sqlx.DB, error) {
	// Build PostgreSQL connection URL.
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	return db.Connect(connURL, "pgx")
}
