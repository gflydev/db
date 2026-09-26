// Package mysql is the MySQL Dialect for github.com/gflydev/db/migrate.
//
// MySQL's DDL statements auto-commit individually — CreateMigrationsTableSQL,
// UpsertMigrationSQL and every migration file run without the rollback safety net that
// PostgreSQL gets from wrapping a file in a transaction. See migrate.Dialect's
// SupportsTransactionalDDL and the spec's "Per-file execution and rollback-on-error" decision
// (docs/specs/2026-09-26-db-migrate-cli.md in the dancefitvn repo).
package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gflydev/core/utils"
	migrate "github.com/gflydev/db/migrate"

	_ "github.com/go-sql-driver/mysql"
)

// Dialect implements migrate.Dialect for MySQL.
type Dialect struct{}

var _ migrate.Dialect = Dialect{}

func (Dialect) Name() string { return "mysql" }

// Open reads the same DB_* environment variables as github.com/gflydev/db/mysql, forcing
// multiStatements (a migration file may contain more than one SQL statement) and parseTime
// (so TIMESTAMP columns scan into time.Time).
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

func (Dialect) UpsertMigrationSQL(table string) string {
	return fmt.Sprintf(`INSERT INTO %s (migration, batch, status, run_count, checksum, migrated_at)
		VALUES (?, ?, 'up', ?, ?, CURRENT_TIMESTAMP)
		ON DUPLICATE KEY UPDATE
			batch = VALUES(batch), status = 'up', run_count = VALUES(run_count),
			checksum = VALUES(checksum), migrated_at = CURRENT_TIMESTAMP`, table)
}

func (Dialect) Placeholder(argPos int) string { return "?" }

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

func (Dialect) Unlock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", key)
	return err
}

func (Dialect) SupportsTransactionalDDL() bool { return false }
