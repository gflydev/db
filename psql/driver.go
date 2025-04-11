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

// ====================================================================
//                           PostgreSQL Driver
// ====================================================================

// New initializes a new PostgreSQL driver and registers it to the database manager.
//
// Returns:
// - *PostgreSQL: A new instance of the PostgreSQL driver.
func New() *PostgreSQL {
	// Set the database type to PostgreSQL in fluentsql.
	fluentsql.SetDBType(fluentsql.PostgreSQL)

	// Create and return a new PostgreSQL driver instance.
	return &PostgreSQL{}
}

// PostgreSQL implements the IDatabase interface for PostgreSQL database operations.
type PostgreSQL struct{}

// Load establishes a connection to the PostgreSQL database.
//
// Returns:
// - *sqlx.DB: The database connection instance.
// - error: An error if the connection fails.
func (d *PostgreSQL) Load() (*sqlx.DB, error) {
	// Build the PostgreSQL connection URL using environment variables.
	// Connection URL format:
	// postgres://username:password@host:port/dbname?sslmode=sslmode
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"), // Database username
		os.Getenv("DB_PASSWORD"), // Database password
		os.Getenv("DB_HOST"),     // Database host
		os.Getenv("DB_PORT"),     // Database port
		os.Getenv("DB_NAME"),     // Database name
		os.Getenv("DB_SSL_MODE"), // SSL mode (e.g., "disable", "require")
	)

	// Establish the database connection using the constructed URL and "pgx" driver.
	return db.Connect(connURL, "pgx")
}
