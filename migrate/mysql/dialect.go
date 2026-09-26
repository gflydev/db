// Package mysql is the MySQL Dialect for github.com/gflydev/db/migrate.
//
// MySQL's DDL statements auto-commit individually — CreateMigrationsTableSQL,
// UpsertMigrationSQL and every migration file run without the rollback safety net that
// PostgreSQL gets from wrapping a file in a transaction. See SupportsTransactionalDDL and the
// migrate package's own migrate.go doc comment on Migrator.runFile for what this means for a
// failure partway through a file.
package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gflydev/core/utils"
	migrate "github.com/gflydev/db/migrate"

	_ "github.com/go-sql-driver/mysql"
)

// Dialect implements migrate.Dialect for MySQL. It is stateless; the zero value (Dialect{}) is
// ready to use.
type Dialect struct{}

var _ migrate.Dialect = Dialect{}

// Name identifies this dialect as "mysql" in log and error output.
func (Dialect) Name() string { return "mysql" }

// Open reads the same DB_HOST/DB_PORT/DB_NAME/DB_USERNAME/DB_PASSWORD environment variables as
// github.com/gflydev/db/mysql and returns a connection pool. It forces multiStatements=true (a
// migration file may contain more than one SQL statement) and parseTime=true (so TIMESTAMP
// columns scan into time.Time) in the DSN regardless of what the caller's environment sets.
func (Dialect) Open() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s?multiStatements=true&parseTime=true",
		utils.Getenv("DB_USERNAME", "user"),
		utils.Getenv("DB_PASSWORD", "secret"),
		utils.Getenv("DB_HOST", "localhost"),
		utils.Getenv("DB_PORT", 3306),
		utils.Getenv("DB_NAME", "gfly"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql: opening connection: %w", err)
	}
	return db, nil
}

// CreateMigrationsTableSQL returns the CREATE TABLE IF NOT EXISTS statement for table, matching
// the schema documented in the migrate package's README (migration, batch, status, run_count,
// checksum, migrated_at, rolled_back_at).
func (Dialect) CreateMigrationsTableSQL(table string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT AUTO_INCREMENT PRIMARY KEY,
		migration VARCHAR(255) NOT NULL UNIQUE,
		batch INT NOT NULL,
		status VARCHAR(10) NOT NULL,
		run_count INT NOT NULL DEFAULT 0,
		checksum VARCHAR(64) NOT NULL,
		migrated_at TIMESTAMP NULL,
		rolled_back_at TIMESTAMP NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`, table)
}

// UpsertMigrationSQL returns an INSERT ... ON DUPLICATE KEY UPDATE statement for table, taking
// (migration, batch, run_count, checksum) as its four "?" placeholders in order and always
// setting status to migrate.StatusUp and migrated_at to the current time.
func (Dialect) UpsertMigrationSQL(table string) string {
	return fmt.Sprintf(`INSERT INTO %s (migration, batch, status, run_count, checksum, migrated_at)
		VALUES (?, ?, '%s', ?, ?, CURRENT_TIMESTAMP)
		ON DUPLICATE KEY UPDATE
			batch = VALUES(batch), status = '%s', run_count = VALUES(run_count),
			checksum = VALUES(checksum), migrated_at = CURRENT_TIMESTAMP`, table, migrate.StatusUp, migrate.StatusUp)
}

// Placeholder always returns "?" — MySQL does not number its placeholders, so argPos is unused.
func (Dialect) Placeholder(argPos int) string { return "?" }

// Lock acquires a connection-scoped MySQL named lock via GET_LOCK, waiting up to 10 seconds.
// Returns an error if the lock could not be acquired within that timeout — most likely because
// another db:migrate process is currently running.
func (Dialect) Lock(ctx context.Context, conn *sql.Conn, key string) error {
	var result sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, 10)", key).Scan(&result); err != nil {
		return err
	}
	if !result.Valid || result.Int64 != 1 {
		return fmt.Errorf("mysql: could not acquire lock %q within 10s — another migration is running", key)
	}
	return nil
}

// Unlock releases the named lock Lock acquired on conn for key.
func (Dialect) Unlock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", key)
	return err
}

// SupportsTransactionalDDL always returns false: MySQL commits each DDL statement immediately,
// so a migration file cannot be rolled back partway through.
func (Dialect) SupportsTransactionalDDL() bool { return false }
