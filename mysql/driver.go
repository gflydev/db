package mysql

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db"
	qb "github.com/jivegroup/fluentsql"
	"github.com/jmoiron/sqlx"

	// Autoload driver for PostgreSQL
	_ "github.com/go-sql-driver/mysql"
)

// ====================================================================
//                           MySQL Driver
// ====================================================================

// New initializes a new MySQL driver instance and registers it to the database manager.
//
// Returns:
//
//	*MySQL: A new instance of the MySQL driver.
func New() *MySQL {
	// Set the database type to MySQL for fluentsql
	qb.SetDialect(new(qb.MySQLDialect))

	// Create and return a new MySQL driver instance
	return &MySQL{}
}

// MySQL is an implementation of the IDatabase interface for MySQL.
type MySQL struct{}

// Load establishes a connection to the MySQL database.
//
// Returns:
//
//	*sqlx.DB: A pointer to the connected database instance.
//	error: An error instance if the connection fails; otherwise, nil.
func (d *MySQL) Load() (*sqlx.DB, error) {
	// Build MySQL connection URL using environment variables or defaults.
	// connURL is a formatted string containing the database connection information.
	connURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s",
		utils.Getenv("DB_USERNAME", "user"),   // Database username
		utils.Getenv("DB_PASSWORD", "secret"), // Database password
		utils.Getenv("DB_HOST", "localhost"),  // Host address
		utils.Getenv("DB_PORT", 3306),         // Port number
		utils.Getenv("DB_NAME", "gfly"),       // Database name
	)

	// Attempt to connect to the database using the constructed connection URL.
	return db.Connect(connURL, "mysql")
}
