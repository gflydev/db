package db

import (
	"github.com/gflydev/core/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

// ========================================================================================
//                                 Structure & Interface
// ========================================================================================

type IDatabase interface {
	Load() (*sqlx.DB, error)
}

// DB the database
type DB struct {
	*sqlx.DB // Embed sqlx DB.
}

// Connect create a connection to Database.
func Connect(connURL, driver string) (*sqlx.DB, error) {
	// Define database connection.
	dbConnection, err := sqlx.Connect(driver, connURL)
	if err != nil {
		return nil, err
	}

	maxConn := utils.Getenv("DB_MAX_CONNECTION", 0)                  // the default is 0 (unlimited)
	maxIdleConn := utils.Getenv("DB_MAX_IDLE_CONNECTION", 2)         // default is 2
	maxLifetimeConn := utils.Getenv("DB_MAX_LIFETIME_CONNECTION", 0) // default is 0, connections are reused forever

	// Set database connection settings:
	dbConnection.SetMaxOpenConns(maxConn)
	dbConnection.SetMaxIdleConns(maxIdleConn)
	dbConnection.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := dbConnection.Ping(); err != nil {
		defer func(db *sqlx.DB) {
			_ = db.Close()
		}(dbConnection)
		return nil, err
	}

	return dbConnection, nil
}

// ========================================================================================
//                                         Drivers
// ========================================================================================

// emptyDB a implement of interface IDatabase for EmptySQL
type emptyDB struct{}

// Load perform DB connection to PostgreSQL database.
func (db *emptyDB) Load() (*sqlx.DB, error) {
	return Connect("empty", "empty")
}

// dbDriver a singleton database driver instance
var dbDriver IDatabase = &emptyDB{}

// Register assign DB provider type fluentsql.PostgreSQL, fluentsql.MySQL,...
func Register(driver IDatabase) {
	dbDriver = driver
}

// ========================================================================================
//                                         Database
// ========================================================================================

// dbInstance a singleton database instance
var dbInstance = &DB{}

// DBInstanceTx returns db transaction instance to handle CRUD at somewhere.
func dbInstanceTx() *sqlx.Tx {
	return dbInstance.MustBegin()
}

// Load func for opening database connection.
func Load() {
	var err error

	dbInstance.DB, err = dbDriver.Load()
	if err != nil {
		panic(err)
	}
}
