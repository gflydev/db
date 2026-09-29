package psql

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db"
	qb "github.com/jivegroup/fluentsql"
	"github.com/jmoiron/sqlx"
	"net"
	"net/url"
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
	// net/url percent-encodes reserved characters (e.g. '@' in a password) so they are not
	// mis-parsed as part of the host.
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			fmt.Sprint(utils.Getenv("DB_USERNAME", "user")),   // Database username
			fmt.Sprint(utils.Getenv("DB_PASSWORD", "secret")), // Database password
		),
		Host: net.JoinHostPort(
			fmt.Sprint(utils.Getenv("DB_HOST", "localhost")), // Host address
			fmt.Sprint(utils.Getenv("DB_PORT", 5432)),        // Port number
		),
		Path:     "/" + fmt.Sprint(utils.Getenv("DB_NAME", "gfly")),                                    // Database name
		RawQuery: url.Values{"sslmode": {fmt.Sprint(utils.Getenv("DB_SSL_MODE", "disable"))}}.Encode(), // SSL mode (e.g., "disable", "require")
	}
	connURL := u.String()

	// Establish the database connection using the constructed URL and "pgx" driver.
	return db.Connect(connURL, "pgx")
}
