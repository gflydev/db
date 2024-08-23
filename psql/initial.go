package psql

import (
	"github.com/gflydev/db"
	_ "github.com/jackc/pgx/v5/stdlib" // Autoload driver for PostgreSQL
	"github.com/jiveio/fluentsql"
)

// ========================================================================================
//                                         Initial
// ========================================================================================

// Auto initial PostgreSQL driver and register to database manager
func init() {
	// Set DBType
	fluentsql.SetDBType(fluentsql.PostgreSQL)

	// Create driver
	provider := &postgreSQL{}

	// Register DB driver
	db.Register(provider)
}
