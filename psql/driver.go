package psql

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db/v2"
	qb "github.com/jivegroup/fluentsql"
	"github.com/jmoiron/sqlx"
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
	// Set the database type to PostgreSQL in qb.
	qb.SetDialect(new(qb.PostgreSQLDialect))

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
	// postgres://username:password@host:port/dbname?sslmode=disable
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		utils.Getenv("DB_USERNAME", "user"),    // Database username
		utils.Getenv("DB_PASSWORD", "secret"),  // Database password
		utils.Getenv("DB_HOST", "localhost"),   // Host address
		utils.Getenv("DB_PORT", 5432),          // Port number
		utils.Getenv("DB_NAME", "gfly"),        // Database name
		utils.Getenv("DB_SSL_MODE", "disable"), // SSL mode (e.g., "disable", "require")
	)

	// Establish the database connection using the constructed URL and "pgx" driver.
	return db.Connect(connURL, "pgx")
}
