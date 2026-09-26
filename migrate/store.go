package migrate

import (
	"context"
	"database/sql"
	"fmt"
)

// store wraps the raw SQL needed to read/write the tracking table, delegating every
// dialect-specific statement to Dialect.
type store struct {
	db      *sql.DB
	dialect Dialect
	table   string
}

func newStore(db *sql.DB, dialect Dialect, table string) *store {
	return &store{db: db, dialect: dialect, table: table}
}

func (s *store) ensureTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, s.dialect.CreateMigrationsTableSQL(s.table))
	if err != nil {
		return fmt.Errorf("migrate: creating %s table: %w", s.table, err)
	}
	return nil
}

// list returns every row currently in the tracking table, keyed by migration name.
func (s *store) list(ctx context.Context) (map[string]Record, error) {
	query := fmt.Sprintf(
		`SELECT migration, batch, status, run_count, checksum, migrated_at, rolled_back_at FROM %s`,
		s.table,
	)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("migrate: listing %s: %w", s.table, err)
	}
	defer rows.Close()

	result := map[string]Record{}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Migration, &r.Batch, &r.Status, &r.RunCount, &r.Checksum, &r.MigratedAt, &r.RolledBackAt); err != nil {
			return nil, fmt.Errorf("migrate: scanning %s row: %w", s.table, err)
		}
		result[r.Migration] = r
	}
	return result, rows.Err()
}

func (s *store) latestBatch(ctx context.Context) (int, error) {
	query := fmt.Sprintf(`SELECT COALESCE(MAX(batch), 0) FROM %s WHERE status = 'up'`, s.table)
	var batch int
	if err := s.db.QueryRowContext(ctx, query).Scan(&batch); err != nil {
		return 0, fmt.Errorf("migrate: reading latest batch from %s: %w", s.table, err)
	}
	return batch, nil
}

// markUp records a fresh first-ever application: run_count is always 1. Re-applications (a
// migration that was rolled back and is now being re-applied) go through markUpWithRunCount so
// the caller (Migrator, which already read the existing row via list) can supply the
// incremented count.
func (s *store) markUp(ctx context.Context, name string, batch int, checksum string) error {
	return s.markUpWithRunCount(ctx, name, batch, checksum, 1)
}

func (s *store) markUpWithRunCount(ctx context.Context, name string, batch int, checksum string, runCount int) error {
	query := s.dialect.UpsertMigrationSQL(s.table)
	_, err := s.db.ExecContext(ctx, query, name, batch, runCount, checksum)
	if err != nil {
		return fmt.Errorf("migrate: marking %s up: %w", name, err)
	}
	return nil
}

func (s *store) markDown(ctx context.Context, name string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET status = 'down', rolled_back_at = CURRENT_TIMESTAMP WHERE migration = %s`,
		s.table, s.dialect.Placeholder(1),
	)
	_, err := s.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("migrate: marking %s down: %w", name, err)
	}
	return nil
}

// baselineMark marks name as applied without executing any SQL file; see Migrator.Baseline for
// the non-empty-table guard this method itself does not enforce.
func (s *store) baselineMark(ctx context.Context, name string, checksum string) error {
	return s.markUpWithRunCount(ctx, name, 0, checksum, 1)
}

func (s *store) isEmpty(ctx context.Context) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.table)
	var count int
	if err := s.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return false, fmt.Errorf("migrate: counting %s: %w", s.table, err)
	}
	return count == 0, nil
}
