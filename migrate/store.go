package migrate

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gflydev/core/errors"
)

// store is the tracking-table data-access layer: every method issues exactly one query or exec
// against table, building dialect-specific SQL only through dialect. Migrator is the only
// caller; store itself has no locking or business-rule logic (batch numbering, checksum
// verification, rollback selection all live in Migrator).
type store struct {
	db      *sql.DB
	dialect Dialect
	table   string
}

// newStore returns a store that reads/writes table through db, using dialect for every
// dialect-specific SQL fragment (the CREATE TABLE statement, the upsert statement, and
// positional placeholders).
func newStore(db *sql.DB, dialect Dialect, table string) *store {
	return &store{db: db, dialect: dialect, table: table}
}

// ensureTable creates the tracking table if it does not already exist. It is safe to call on
// every Migrator operation — Migrator.prepare does exactly that — since the dialect's SQL uses
// CREATE TABLE IF NOT EXISTS.
func (s *store) ensureTable(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, s.dialect.CreateMigrationsTableSQL(s.table)); err != nil {
		return errors.New("creating %s table: %w", s.table, err)
	}
	return nil
}

// list returns every row currently in the tracking table, keyed by migration name. The caller
// (Migrator) joins this against the files Load found on disk; a name with no entry here has
// never been applied.
func (s *store) list(ctx context.Context) (map[string]Record, error) {
	query := fmt.Sprintf(
		`SELECT migration, batch, status, run_count, checksum, migrated_at, rolled_back_at FROM %s`,
		s.table,
	)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("listing %s: %w", s.table, err)
	}
	defer rows.Close()

	result := map[string]Record{}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Migration, &r.Batch, &r.Status, &r.RunCount, &r.Checksum, &r.MigratedAt, &r.RolledBackAt); err != nil {
			return nil, errors.New("scanning %s row: %w", s.table, err)
		}
		result[r.Migration] = r
	}
	return result, rows.Err()
}

// latestBatch returns the highest batch number among rows currently marked StatusUp, or 0 if
// none are. Migrator.Up increments this to number a fresh batch.
func (s *store) latestBatch(ctx context.Context) (int, error) {
	query := fmt.Sprintf(`SELECT COALESCE(MAX(batch), 0) FROM %s WHERE status = '%s'`, s.table, StatusUp)
	var batch int
	if err := s.db.QueryRowContext(ctx, query).Scan(&batch); err != nil {
		return 0, errors.New("reading latest batch from %s: %w", s.table, err)
	}
	return batch, nil
}

// markUp upserts name's row as StatusUp with the given batch, checksum and runCount, and sets
// migrated_at to now. It is used both for a fresh first-ever application (runCount == 1, called
// by Migrator.Up) and for a re-application after a rollback (runCount == the previous RunCount
// + 1) — the caller is responsible for computing runCount from the row store.list already
// returned, since store has no memory of prior state beyond what's in the table.
func (s *store) markUp(ctx context.Context, name string, batch int, checksum string, runCount int) error {
	query := s.dialect.UpsertMigrationSQL(s.table)
	if _, err := s.db.ExecContext(ctx, query, name, batch, runCount, checksum); err != nil {
		return errors.New("marking %s up: %w", name, err)
	}
	return nil
}

// markDown sets name's row to StatusDown and stamps rolled_back_at. It leaves batch and
// run_count untouched, so they remain a historical record of the last time this migration was
// applied (see Record.Batch).
func (s *store) markDown(ctx context.Context, name string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET status = '%s', rolled_back_at = CURRENT_TIMESTAMP WHERE migration = %s`,
		s.table, StatusDown, s.dialect.Placeholder(1),
	)
	if _, err := s.db.ExecContext(ctx, query, name); err != nil {
		return errors.New("marking %s down: %w", name, err)
	}
	return nil
}

// baselineMark upserts name's row as StatusUp in batch 0 with run_count 1, exactly like a fresh
// markUp call, but exists as its own method so Migrator.Baseline's intent (bootstrap, not a real
// application) is visible at the call site. See Migrator.Baseline for the non-empty-table guard
// this method itself does not enforce.
func (s *store) baselineMark(ctx context.Context, name string, checksum string) error {
	return s.markUp(ctx, name, 0, checksum, 1)
}

// isEmpty reports whether the tracking table currently has zero rows. Migrator.Baseline uses
// this to refuse running against a table that already has real history, unless force is true.
func (s *store) isEmpty(ctx context.Context) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.table)
	var count int
	if err := s.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return false, errors.New("counting %s: %w", s.table, err)
	}
	return count == 0, nil
}
