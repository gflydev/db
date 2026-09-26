// Package migrate runs versioned *.sql migration files against a database, tracking each
// file's state (applied/rolled back, batch, run count, checksum) in a tracking table it manages
// itself. It is dialect-agnostic: see the postgres and mysql subpackages for the two supported
// databases, and RunCLI for wiring it into an application's own command-line entrypoint.
//
// # File layout
//
// Migrations live as pairs of files named "NNNNNN_description.up.sql" and
// "NNNNNN_description.down.sql" (golang-migrate's convention) in a single directory; Load
// discovers and validates them.
//
// # Concurrency
//
// A Migrator is not safe for concurrent use by multiple goroutines within one process — each
// exported method holds the underlying *sql.DB open for its own duration and returns before the
// next call should start. Concurrent processes are handled separately: every method acquires
// the dialect's advisory lock for its full duration, so two OS processes (e.g. two deploys
// racing) serialize instead of corrupting the tracking table.
package migrate

import (
	"context"
	"database/sql"
)

const (
	// DefaultTable is the tracking table name every adopter of this package uses. It is not
	// configurable (see NewMigrator) — one gFly app has one migrations directory and one
	// tracking table, and a configurable name would only invite drift between the table a
	// health check reads and the one db:migrate writes.
	DefaultTable = "migrations"

	// StatusUp and StatusDown are the two values Record.Status takes. They are exported so a
	// caller inspecting Status's results doesn't need to guess the tracking table's raw string
	// values, and so this package's own SQL (store.go) and each Dialect's generated SQL
	// (postgres, mysql) share one definition instead of three independently hand-typed copies.
	StatusUp   = "up"
	StatusDown = "down"
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

// Record is one row of the tracking table, as last read by store.list.
type Record struct {
	// Migration is the tracked migration's identity — the same value as the matching
	// Migration.Name (e.g. "000024_create_widgets_table").
	Migration string
	// Batch is the batch number assigned the last time this migration was applied. It keeps
	// its value after a rollback (Status becomes StatusDown) so it remains a historical record;
	// it is only reassigned by a fresh "up".
	Batch int
	// Status is either StatusUp or StatusDown.
	Status string
	// RunCount is the number of times this migration has been successfully applied, including
	// re-applications after a rollback. It is never decremented.
	RunCount int
	// Checksum is the sha256 (hex-encoded) of the .up.sql file's contents as of the last time
	// this migration was applied. It is compared against the current file on disk before every
	// operation; see ErrChecksumMismatch.
	Checksum string
	// MigratedAt is the time of the last successful "up"; zero-valued (Valid == false) if this
	// migration has never been applied.
	MigratedAt sql.NullTime
	// RolledBackAt is the time of the last "down"; zero-valued (Valid == false) if this
	// migration has never been rolled back.
	RolledBackAt sql.NullTime
}

// StatusRow is one line of `db:migrate --status` output: a Migration joined with its Record,
// if any, plus whether its on-disk checksum still matches what was recorded.
type StatusRow struct {
	Migration       Migration
	Record          *Record // nil if never applied
	ChecksumMatches bool    // meaningless (true) if Record is nil
}
