// Package postgres is the PostgreSQL Dialect for github.com/gflydev/db/migrate.
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

// Dialect implements migrate.Dialect for PostgreSQL.
type Dialect struct{}

var _ migrate.Dialect = Dialect{}

func (Dialect) Name() string { return "postgres" }

// Open reads the same DB_* environment variables as github.com/gflydev/db/psql.
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

func (Dialect) UpsertMigrationSQL(table string) string {
	return fmt.Sprintf(`INSERT INTO %s (migration, batch, status, run_count, checksum, migrated_at)
		VALUES ($1, $2, 'up', $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (migration) DO UPDATE SET
			batch = EXCLUDED.batch, status = 'up', run_count = EXCLUDED.run_count,
			checksum = EXCLUDED.checksum, migrated_at = CURRENT_TIMESTAMP`, table)
}

func (Dialect) Placeholder(argPos int) string {
	return fmt.Sprintf("$%d", argPos)
}

func (Dialect) Lock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock(hashtext($1))", key)
	return err
}

func (Dialect) Unlock(ctx context.Context, conn *sql.Conn, key string) error {
	_, err := conn.ExecContext(ctx, "SELECT pg_advisory_unlock(hashtext($1))", key)
	return err
}

func (Dialect) SupportsTransactionalDDL() bool { return true }
