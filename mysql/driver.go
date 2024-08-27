package mysql

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db"
	"github.com/jiveio/fluentsql"
	"github.com/jmoiron/sqlx"

	// Autoload driver for PostgreSQL
	_ "github.com/go-sql-driver/mysql"
)

// ========================================================================================
//                                     MySQL Driver
// ========================================================================================

// New initial MySQL driver and register to database manager
func New() *MySQL {
	// Set DBType
	fluentsql.SetDBType(fluentsql.MySQL)

	// Create driver
	return &MySQL{}
}

// MySQL a implement of interface IDatabase for MySQL
type MySQL struct{}

// Load perform DB connection to Mysql database.
func (d *MySQL) Load() (*sqlx.DB, error) {
	// Build Mysql connection URL.
	connURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s",
		utils.Getenv("DB_USERNAME", "root"),
		utils.Getenv("DB_PASSWORD", "secret"),
		utils.Getenv("DB_HOST", "localhost"),
		utils.Getenv("DB_PORT", 3306),
		utils.Getenv("DB_NAME", "gfly"),
	)

	return db.Connect(connURL, "mysql")
}
