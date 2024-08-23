package mysql

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/db"
	"github.com/jmoiron/sqlx"
)

// ========================================================================================
//                                     MySQL Driver
// ========================================================================================

// mySQL a implement of interface IDatabase for MySQL
type mySQL struct{}

// Load perform DB connection to Mysql database.
func (d *mySQL) Load() (*sqlx.DB, error) {
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
