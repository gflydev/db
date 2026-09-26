// Package migrate runs versioned *.sql migration files against a database, tracking each
// file's state (applied/rolled back, batch, run count, checksum) in a "migrations" table it
// manages itself. It is dialect-agnostic: see the postgres and mysql subpackages for the two
// supported databases.
package migrate

import (
	"context"
	"database/sql"
)

// Dialect is the seam between this package's orchestration and a specific database.
type Dialect interface {
	// Name identifies the dialect in log/error output, e.g. "postgres", "mysql".
	Name() string

	// Open reads DB_HOST/DB_PORT/DB_NAME/DB_USERNAME/DB_PASSWORD/DB_SSL_MODE and returns a
	// ready connection pool.
	Open() (*sql.DB, error)

	// CreateMigrationsTableSQL returns the CREATE TABLE IF NOT EXISTS statement for the
	// tracking table named table.
	CreateMigrationsTableSQL(table string) string

	// UpsertMigrationSQL returns a statement that inserts a new row for table, or updates the
	// existing row for the same migration name, setting status='up'. Parameter order is
	// (migration, batch, run_count, checksum).
	UpsertMigrationSQL(table string) string

	// Placeholder returns the positional-parameter placeholder for argument position argPos
	// (1-based): "$1", "$2", ... for Postgres, "?" for every position in MySQL.
	Placeholder(argPos int) string

	// Lock acquires a database-wide advisory lock scoped to conn, blocking (up to the
	// dialect's own timeout) until it is available or returning an error if it cannot be
	// acquired.
	Lock(ctx context.Context, conn *sql.Conn, key string) error

	// Unlock releases a lock acquired by Lock, using the same conn.
	Unlock(ctx context.Context, conn *sql.Conn, key string) error

	// SupportsTransactionalDDL reports whether a failing statement partway through a
	// migration file can be rolled back (true for Postgres, false for MySQL).
	SupportsTransactionalDDL() bool
}

// Migration is one file pair discovered on disk by Load.
type Migration struct {
	// Version is the 6-digit numeric prefix, e.g. "000024".
	Version string
	// Name is the filename's identity: the numeric prefix and description, without the
	// .up.sql/.down.sql suffix, e.g. "000024_create_widgets_table".
	Name string
	// UpPath and DownPath are absolute (or working-directory-relative) paths to the two files.
	UpPath, DownPath string
}

// Record is one row of the migrations table.
type Record struct {
	Migration    string
	Batch        int
	Status       string // "up" or "down"
	RunCount     int
	Checksum     string
	MigratedAt   sql.NullTime
	RolledBackAt sql.NullTime
}

// StatusRow is one line of `db:migrate --status` output: a Migration joined with its Record,
// if any, plus whether its on-disk checksum still matches what was recorded.
type StatusRow struct {
	Migration       Migration
	Record          *Record // nil if never applied
	ChecksumMatches bool    // meaningless (true) if Record is nil
}
