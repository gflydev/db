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
// This interface provides a standardized way to implement different database drivers
// and connection strategies while maintaining a consistent API for database initialization.
// Implementations should handle driver-specific connection logic and configuration.
//
// Methods:
//   - Load() (*sqlx.DB, error): Establishes and returns a configured database connection
//
// Examples of implementations:
//   - MySQL driver: Handles MySQL-specific connection strings and configurations
//   - PostgreSQL driver: Handles PostgreSQL-specific connection strings and configurations
//   - Mock/Test driver: Provides test database connections for unit testing
//
// Usage:
//
//	type MySQLDriver struct {
//	    connectionString string
//	}
//
//	func (d *MySQLDriver) Load() (*sqlx.DB, error) {
//	    return Connect(d.connectionString, "mysql")
//	}
//
//	// Register the driver
//	Register(&MySQLDriver{connectionString: "user:pass@tcp(localhost:3306)/dbname"})
type IDatabase interface {
	// Load establishes a connection to the database and returns a configured sqlx.DB instance.
	// This method should handle all driver-specific connection logic, including:
	// - Connection string parsing and validation
	// - Driver-specific configuration options
	// - Connection pool settings
	// - Initial connection verification
	//
	// Returns:
	//   - *sqlx.DB: A fully configured and tested database connection ready for use
	//   - error: An error if connection establishment fails, including network issues,
	//     authentication failures, or configuration problems
	Load() (*sqlx.DB, error)
}

// DB represents a database connection wrapper that embeds sqlx.DB to provide extended functionality.
// This struct serves as the foundation for database operations throughout the application,
// providing a consistent interface while leveraging the powerful features of sqlx.DB.
// The embedded sqlx.DB allows direct access to all sqlx methods while enabling
// additional custom functionality to be added as needed.
//
// Fields:
//   - *sqlx.DB: Embedded sqlx.DB instance providing:
//   - Named parameter support for queries
//   - Struct scanning capabilities
//   - Transaction management
//   - Connection pooling
//   - Prepared statement caching
//
// Usage:
//
//	// The DB struct is typically used internally by the ORM
//	// and accessed through the global dbInstance variable
//	db := &DB{DB: sqlxConnection}
//
//	// Direct access to sqlx methods is available
//	rows, err := dbInstance.Query("SELECT * FROM users")
//
//	// Custom methods can be added to extend functionality
//	// while maintaining compatibility with sqlx.DB
type DB struct {
	*sqlx.DB // Embedded sqlx.DB for working with SQL databases, providing full sqlx functionality
}

// Connect establishes a database connection with comprehensive configuration and connection pooling.
// This function creates a new database connection using the specified driver and connection URL,
// then applies optimal connection pool settings based on environment variables or defaults.
// It also performs connection validation through a ping operation to ensure the database
// is accessible and responsive before returning the connection instance.
//
// Parameters:
//   - connURL (string): The database connection URL/DSN (Data Source Name). Format varies by driver:
//   - MySQL: "user:password@tcp(host:port)/database?param=value"
//   - PostgreSQL: "postgres://user:password@host:port/database?sslmode=disable"
//   - SQLite: "file:path/to/database.db" or ":memory:" for in-memory database
//   - driver (string): The database driver name. Supported values:
//   - "mysql": MySQL database driver
//   - "postgres" or "postgresql": PostgreSQL database driver
//   - "sqlite3": SQLite database driver
//   - Custom drivers registered with database/sql
//
// Returns:
//   - *sqlx.DB: A fully configured database connection with optimized pool settings:
//   - Connection pooling configured based on environment variables
//   - Connection lifetime and idle time management
//   - Verified connectivity through ping operation
//   - error: Returns an error if:
//   - Driver is not registered or invalid
//   - Connection URL is malformed or invalid
//   - Database server is unreachable
//   - Authentication fails
//   - Database does not exist
//   - Network connectivity issues occur
//
// Environment Variables (with defaults):
//   - DB_MAX_CONNECTION: Maximum open connections (default: 0 = unlimited)
//   - DB_MAX_IDLE_CONNECTION: Maximum idle connections (default: 10)
//   - DB_MAX_LIFETIME_CONNECTION: Connection lifetime in minutes (default: 30)
//   - DB_MAX_IDLE_TIME_CONNECTION: Connection idle time in minutes (default: 3)
//
// Examples:
//
//	// MySQL connection
//	db, err := Connect("user:pass@tcp(localhost:3306)/mydb", "mysql")
//	if err != nil {
//	    log.Fatal("Failed to connect to MySQL:", err)
//	}
//	defer db.Close()
//
//	// PostgreSQL connection
//	connURL := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
//	db, err := Connect(connURL, "postgres")
//	if err != nil {
//	    log.Fatal("Failed to connect to PostgreSQL:", err)
//	}
//
//	// SQLite connection
//	db, err := Connect("file:./app.db", "sqlite3")
//	if err != nil {
//	    log.Fatal("Failed to connect to SQLite:", err)
//	}
//
//	// In-memory SQLite for testing
//	db, err := Connect(":memory:", "sqlite3")
//	if err != nil {
//	    log.Fatal("Failed to create in-memory database:", err)
//	}
//
// Note:
//   - Connection pool settings are automatically applied for optimal performance
//   - The connection is validated with a ping operation before being returned
//   - Failed connections are automatically closed to prevent resource leaks
//   - Environment variables allow runtime configuration without code changes
func Connect(connURL, driver string) (*sqlx.DB, error) {
	// Define database connection.
	dbConnection, err := sqlx.Connect(driver, connURL)
	if err != nil {
		return nil, err
	}

	// Load configuration settings for database connections from environment variables.
	maxConn := utils.Getenv("DB_MAX_CONNECTION", 0)                   // Maximum open connections (default: 0, unlimited).
	maxIdleConn := utils.Getenv("DB_MAX_IDLE_CONNECTION", 10)         // Maximum idle connections (default: 10).
	maxLifetimeConn := utils.Getenv("DB_MAX_LIFETIME_CONNECTION", 30) // Maximum lifetime of connections in minutes (default: 30).
	maxIdleTimeConn := utils.Getenv("DB_MAX_IDLE_TIME_CONNECTION", 3) // Maximum idle time of connections in minutes (default: 3).

	// Set database connection settings.
	dbConnection.SetMaxOpenConns(maxConn)
	dbConnection.SetMaxIdleConns(maxIdleConn)
	dbConnection.SetConnMaxLifetime(time.Duration(maxLifetimeConn) * time.Minute)
	dbConnection.SetConnMaxIdleTime(time.Duration(maxIdleTimeConn) * time.Minute)

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

// emptyDB implements the IDatabase interface as a default no-operation database driver.
// This struct serves as a placeholder driver that prevents nil pointer errors when
// no actual database driver has been registered. It's designed to fail gracefully
// when attempting to establish a real database connection, making it useful for
// testing scenarios or as a safe default state.
//
// Usage:
//   - Automatically used as the default driver when no custom driver is registered
//   - Provides a safe fallback that won't cause runtime panics
//   - Useful for testing environments where database connections aren't needed
//   - Can be replaced with actual drivers via the Register() function
//
// Example:
//
//	// This will use emptyDB by default
//	Load() // Will attempt to connect with "empty" driver and fail gracefully
//
//	// Register a real driver to replace emptyDB
//	Register(&MySQLDriver{})
//	Load() // Now uses the registered MySQL driver
type emptyDB struct{}

// Load attempts to establish a mock database connection using placeholder values.
// This method is designed to fail gracefully, providing a safe default behavior
// when no real database driver has been registered. It serves as a placeholder
// implementation that prevents nil pointer errors while clearly indicating
// that a proper database driver should be registered for production use.
//
// Returns:
//   - *sqlx.DB: Will be nil since the connection with "empty" driver will fail
//   - error: Always returns an error because "empty" is not a valid database driver.
//     This is intentional behavior to indicate that a proper driver should be registered.
//
// Examples:
//
//	// Direct usage (not recommended for production)
//	emptyDriver := &emptyDB{}
//	db, err := emptyDriver.Load()
//	if err != nil {
//	    log.Println("Expected error: empty driver cannot establish real connections")
//	}
//
//	// Typical usage through the driver system
//	// Before registering a real driver, this will be called automatically
//	Load() // Uses emptyDB.Load() internally and will fail
//
// Note:
//   - This method is intended to fail as a safety mechanism
//   - The failure indicates that a proper database driver should be registered
//   - Used internally as the default driver to prevent nil pointer panics
//   - Should be replaced with actual database drivers in production environments
func (db *emptyDB) Load() (*sqlx.DB, error) {
	return Connect("empty", "empty") // Intentionally fails with invalid driver to prompt proper driver registration
}

// dbDriver holds a singleton instance of the currently registered database driver.
// Defaults to the emptyDB instance.
var dbDriver IDatabase = &emptyDB{}

// Register replaces the default database driver with a custom implementation.
// This function allows applications to specify their database driver (MySQL, PostgreSQL, etc.)
// by providing an implementation of the IDatabase interface. The registered driver will be
// used by the Load() function to establish database connections. This design enables
// driver-specific configuration and connection logic while maintaining a consistent API.
//
// Parameters:
//   - driver (IDatabase): A custom database driver that implements the IDatabase interface.
//     The driver should handle:
//   - Database-specific connection string formatting
//   - Driver-specific configuration options
//   - Connection establishment and validation
//   - Error handling for connection failures
//
// Examples:
//
//	// MySQL driver registration
//	type MySQLDriver struct {
//	    Host     string
//	    Port     int
//	    Database string
//	    Username string
//	    Password string
//	}
//
//	func (d *MySQLDriver) Load() (*sqlx.DB, error) {
//	    connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
//	        d.Username, d.Password, d.Host, d.Port, d.Database)
//	    return Connect(connStr, "mysql")
//	}
//
//	// Register the MySQL driver
//	mysqlDriver := &MySQLDriver{
//	    Host:     "localhost",
//	    Port:     3306,
//	    Database: "myapp",
//	    Username: "user",
//	    Password: "password",
//	}
//	Register(mysqlDriver)
//
//	// PostgreSQL driver registration
//	type PostgreSQLDriver struct {
//	    ConnectionString string
//	}
//
//	func (d *PostgreSQLDriver) Load() (*sqlx.DB, error) {
//	    return Connect(d.ConnectionString, "postgres")
//	}
//
//	pgDriver := &PostgreSQLDriver{
//	    ConnectionString: "postgres://user:pass@localhost/dbname?sslmode=disable",
//	}
//	Register(pgDriver)
//
//	// Test/Mock driver registration
//	type MockDriver struct{}
//
//	func (d *MockDriver) Load() (*sqlx.DB, error) {
//	    return Connect(":memory:", "sqlite3")
//	}
//
//	Register(&MockDriver{}) // Useful for testing
//
// Note:
//   - This function replaces the global dbDriver variable
//   - Should be called before Load() to ensure the correct driver is used
//   - The driver will be used for all subsequent database connections
//   - Typically called during application initialization
//   - Can be called multiple times to change drivers (useful for testing)
func Register(driver IDatabase) {
	dbDriver = driver
}

// ====================================================================
//                              Database
// ====================================================================

// dbInstance is a singleton instance of the DB struct used for managing database operations.
var dbInstance = &DB{}

// dbInstanceTx creates and returns a new database transaction from the global database instance.
// This function provides a convenient way to start database transactions for operations
// that require atomicity, consistency, isolation, and durability (ACID properties).
// It uses MustBegin() which will panic if the transaction cannot be started, ensuring
// that transaction failures are immediately apparent rather than silently ignored.
//
// Returns:
//   - *sqlx.Tx: A new database transaction instance that provides:
//   - Transaction-scoped database operations
//   - Rollback capabilities for error handling
//   - Commit functionality to persist changes
//   - All standard sqlx.Tx methods for querying and execution
//
// Examples:
//
//	// Basic transaction usage
//	tx := dbInstanceTx()
//	defer func() {
//	    if r := recover(); r != nil {
//	        tx.Rollback() // Rollback on panic
//	    }
//	}()
//
//	// Perform transactional operations
//	_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "John")
//	if err != nil {
//	    tx.Rollback()
//	    return err
//	}
//
//	// Commit the transaction
//	return tx.Commit()
//
//	// Transaction with multiple operations
//	tx := dbInstanceTx()
//	defer tx.Rollback() // Rollback if not committed
//
//	// Multiple related operations
//	userID, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "Jane")
//	if err != nil {
//	    return err
//	}
//
//	_, err = tx.Exec("INSERT INTO profiles (user_id, bio) VALUES (?, ?)", userID, "Bio")
//	if err != nil {
//	    return err
//	}
//
//	return tx.Commit() // Commit all changes
//
// Note:
//   - Uses MustBegin() which panics on failure rather than returning an error
//   - Requires the global dbInstance to be properly initialized via Load()
//   - Transactions should always be committed or rolled back to avoid resource leaks
//   - Use defer statements for automatic rollback in error scenarios
//   - Nested transactions are not supported by most databases
func dbInstanceTx() *sqlx.Tx {
	return dbInstance.MustBegin()
}

// Load initializes the global database connection using the registered driver.
// This function establishes the primary database connection that will be used throughout
// the application lifecycle. It delegates the actual connection establishment to the
// registered database driver and assigns the result to the global dbInstance variable.
// The function is designed to be called once during application startup.
//
// Behavior:
//   - Calls the Load() method on the currently registered database driver
//   - Assigns the resulting connection to the global dbInstance variable
//   - Panics immediately if the connection cannot be established
//   - Should be called after registering a proper database driver via Register()
//
// Panics:
//   - If no database driver has been registered (uses emptyDB which always fails)
//   - If the registered driver fails to establish a connection
//   - If database server is unreachable or credentials are invalid
//   - If the specified database does not exist
//
// Examples:
//
//	// Basic application initialization
//	func main() {
//	    // Register database driver
//	    mysqlDriver := &MySQLDriver{
//	        Host: "localhost",
//	        Port: 3306,
//	        Database: "myapp",
//	        Username: "user",
//	        Password: "password",
//	    }
//	    Register(mysqlDriver)
//
//	    // Initialize database connection
//	    Load() // Will panic if connection fails
//
//	    // Application is ready to use database
//	    // dbInstance is now available for ORM operations
//	}
//
//	// With error handling (using recover)
//	func initDatabase() (err error) {
//	    defer func() {
//	        if r := recover(); r != nil {
//	            err = fmt.Errorf("database initialization failed: %v", r)
//	        }
//	    }()
//
//	    Register(&PostgreSQLDriver{ConnectionString: "postgres://..."})
//	    Load()
//	    return nil
//	}
//
//	// Testing setup
//	func setupTestDB() {
//	    Register(&MockDriver{}) // Use in-memory database
//	    Load()
//	}
//
// Note:
//   - Must be called after Register() to use a proper database driver
//   - Panics on failure to ensure database connectivity is verified at startup
//   - Should only be called once during application initialization
//   - The resulting connection is stored in the global dbInstance variable
//   - All ORM operations depend on this function being called successfully
func Load() {
	var err error

	// Load the database connection using the registered driver.
	if dbInstance.DB, err = dbDriver.Load(); err != nil {
		// If an error occurs, panic to prevent further execution.
		panic(err)
	}
}
