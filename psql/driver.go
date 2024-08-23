package psql

import (
	"fmt"
	"github.com/gflydev/db"
	"github.com/jmoiron/sqlx"
	"os"
)

// ========================================================================================
//                                     PostgreSQL Driver
// ========================================================================================

// postgreSQL a implement of interface IDatabase for PostgreSQL
type postgreSQL struct{}

// Load perform DB connection to PostgreSQL database.
func (d *postgreSQL) Load() (*sqlx.DB, error) {
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
