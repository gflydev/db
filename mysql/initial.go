package mysql

import (
	"github.com/gflydev/db"
	_ "github.com/go-sql-driver/mysql" // Autoload driver for Mysql
	"github.com/jiveio/fluentsql"
)

// ========================================================================================
//                                         Initial
// ========================================================================================

// Auto initial MySQL driver and register to database manager
func init() {
	// Set DBType
	fluentsql.SetDBType(fluentsql.MySQL)

	// Create driver
	provider := &mySQL{}

	// Register DB driver
	db.Register(provider)
}
