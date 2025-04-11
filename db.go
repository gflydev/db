package db

import (
	"github.com/gflydev/core/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

// ====================================================================
//                          Structure & Interface
// ====================================================================

// IDatabase interface defines a contract for database loading operations.
type IDatabase interface {
	// Load establishes a connection to the database and returns the instance of sqlx.DB or an error if the connection fails.
	Load() (*sqlx.DB, error)
}

// DB represents a database and embeds the sqlx.DB to provide extended functionality.
type DB struct {
	*sqlx.DB // Embedded sqlx.DB for working with SQL databases.
}

// Connect creates a connection to a database using the provided connection URL and driver name.
// Parameters:
// - connURL: A string representing the connection URL to the database.
// - driver: A string representing the name of the database driver (e.g., "postgres", "mysql").
// Returns:
// - *sqlx.DB: A pointer to the sqlx.DB instance representing the database connection.
// - error: An error if the connection initialization fails.
func Connect(connURL, driver string) (*sqlx.DB, error) {
	// Define database connection.
	dbConnection, err := sqlx.Connect(driver, connURL)
	if err != nil {
		return nil, err
	}

	// Load configuration settings for database connections from environment variables.
	maxConn := utils.Getenv("DB_MAX_CONNECTION", 0)                  // Maximum open connections (default: 0, unlimited).
	maxIdleConn := utils.Getenv("DB_MAX_IDLE_CONNECTION", 2)         // Maximum idle connections (default: 2).
	maxLifetimeConn := utils.Getenv("DB_MAX_LIFETIME_CONNECTION", 0) // Maximum lifetime of connections in nanoseconds (default: 0, infinite).

	// Set database connection settings.
	dbConnection.SetMaxOpenConns(maxConn)
	dbConnection.SetMaxIdleConns(maxIdleConn)
	dbConnection.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database to verify the connection.
	if err := dbConnection.Ping(); err != nil {
		// Close the connection on error.
		defer func(db *sqlx.DB) {
			_ = db.Close()
		}(dbConnection)
		return nil, err
	}

	return dbConnection, nil
}

// ====================================================================
//                              Drivers
// ====================================================================

// emptyDB implements the IDatabase interface for an empty (no-op) SQL driver.
type emptyDB struct{}

// Load establishes a mock connection to a database using the empty driver.
// Returns:
// - *sqlx.DB: A pointer to the sqlx.DB instance representing the database connection.
// - error: An error if the connection initialization fails.
func (db *emptyDB) Load() (*sqlx.DB, error) {
	return Connect("empty", "empty") // Mock connection with empty arguments.
}

// dbDriver holds a singleton instance of the currently registered database driver.
// Defaults to the emptyDB instance.
var dbDriver IDatabase = &emptyDB{}

// Register assigns a custom database driver to the dbDriver variable.
// Parameters:
// - driver (IDatabase): The custom database driver implementing the IDatabase interface.
func Register(driver IDatabase) {
	dbDriver = driver
}

// ====================================================================
//                              Database
// ====================================================================

// dbInstance is a singleton instance of the DB struct used for managing database operations.
var dbInstance = &DB{}

// dbInstanceTx begins a new database transaction.
// Returns:
// - *sqlx.Tx: A pointer to a new sqlx.Tx instance representing the database transaction.
func dbInstanceTx() *sqlx.Tx {
	return dbInstance.MustBegin()
}

// Load initializes the database connection by loading it through the registered database driver and assigns it to the dbInstance.
// Panics if the database connection fails.
func Load() {
	var err error

	// Load the database connection using the registered driver.
	dbInstance.DB, err = dbDriver.Load()
	if err != nil {
		// If an error occurs, panic to prevent further execution.
		panic(err)
	}
}
