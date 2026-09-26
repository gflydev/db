// Package postgres is the PostgreSQL Dialect for github.com/gflydev/db/migrate.
//
// Every migration file runs inside its own transaction (SupportsTransactionalDDL reports true):
// PostgreSQL allows DDL inside a transaction, so a failing statement leaves no partial effect.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gflydev/core/utils"
	migrate "github.com/gflydev/db/migrate"

	// Autoload the PostgreSQL driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Dialect implements migrate.Dialect for PostgreSQL. It is stateless; the zero value
// (Dialect{}) is ready to use.
type Dialect struct{}

var _ migrate.Dialect = Dialect{}

// Name identifies this dialect as "postgres" in log and error output.
func (Dialect) Name() string { return "postgres" }

// Open reads the same DB_HOST/DB_PORT/DB_NAME/DB_USERNAME/DB_PASSWORD/DB_SSL_MODE environment
// variables as github.com/gflydev/db/psql and returns a pgx-backed connection pool. It does not
// verify connectivity — callers that need to fail fast on a bad connection should call
// (*sql.DB).PingContext themselves.
func (Dialect) Open() (*sql.DB, error) {
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		utils.Getenv("DB_USERNAME", "user"),
		utils.Getenv("DB_PASSWORD", "secret"),
		utils.Getenv("DB_HOST", "localhost"),
		utils.Getenv("DB_PORT", 5432),
		utils.Getenv("DB_NAME", "gfly"),
		utils.Getenv("DB_SSL_MODE", "disable"),
	)
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return nil, fmt.Errorf("postgres: opening connection: %w", err)
	}
	return db, nil
}

// CreateMigrationsTableSQL returns the CREATE TABLE IF NOT EXISTS statement for table, matching
// the schema documented in the migrate package's README (migration, batch, status, run_count,
// checksum, migrated_at, rolled_back_at).
func (Dialect) CreateMigrationsTableSQL(table string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		migration VARCHAR(255) NOT NULL UNIQUE,
		batch INTEGER NOT NULL,
		status VARCHAR(10) NOT NULL,
		run_count INTEGER NOT NULL DEFAULT 0,
		checksum VARCHAR(64) NOT NULL,
		migrated_at TIMESTAMP NULL,
		rolled_back_at TIMESTAMP NULL
	)`, table)
}

// UpsertMigrationSQL returns an INSERT ... ON CONFLICT DO UPDATE statement for table, taking
// (migration, batch, run_count, checksum) as $1-$4 and always setting status to migrate.StatusUp
// and migrated_at to the current time.
func (Dialect) UpsertMigrationSQL(table string) string {
	return fmt.Sprintf(`INSERT INTO %s (migration, batch, status, run_count, checksum, migrated_at)
		VALUES ($1, $2, '%s', $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (migration) DO UPDATE SET
			batch = EXCLUDED.batch, status = '%s', run_count = EXCLUDED.run_count,
			checksum = EXCLUDED.checksum, migrated_at = CURRENT_TIMESTAMP`, table, migrate.StatusUp, migrate.StatusUp)
}

// Placeholder returns PostgreSQL's numbered placeholder ("$1", "$2", ...) for argument position
// argPos.
func (Dialect) Placeholder(argPos int) string {
	return fmt.Sprintf("$%d", argPos)
}

// Lock acquires a session-scoped PostgreSQL advisory lock keyed by the hash of key, blocking
// until it is available. The lock is released by Unlock on the same conn — advisory locks are
// tied to the session (connection) that took them, not to a transaction.
func (Dialect) Lock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock(hashtext($1))", key)
	return err
}

// Unlock releases the advisory lock Lock acquired on conn for key.
func (Dialect) Unlock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_unlock(hashtext($1))", key)
	return err
}

// SupportsTransactionalDDL always returns true: PostgreSQL rolls back DDL inside a failed
// transaction like any other statement.
func (Dialect) SupportsTransactionalDDL() bool { return true }
