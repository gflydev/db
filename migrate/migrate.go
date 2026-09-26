package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
)

// Result reports what a single Up/Down invocation did, in the order it happened.
type Result struct {
	Applied []string
}

// Migrator runs *.sql files from dir against db, tracking state in table via dialect.
type Migrator struct {
	db      *sql.DB
	dialect Dialect
	dir     string
	table   string
	store   *store
}

// NewMigrator wires up a Migrator. db is the pool to run against; dialect selects the SQL
// dialect; dir is the migrations directory; table is the tracking table name (every adopter
// uses "migrations" — the spec fixes this rather than making it configurable).
func NewMigrator(db *sql.DB, dialect Dialect, dir string, table string) *Migrator {
	return &Migrator{
		db:      db,
		dialect: dialect,
		dir:     dir,
		table:   table,
		store:   newStore(db, dialect, table),
	}
}

const lockKey = "gfly_db_migrate"

// prepare ensures the tracking table exists, loads the migration files, loads the current
// tracking rows, and verifies every "up" migration's on-disk checksum still matches what was
// recorded. Every public method calls this first (Status passes force=true so it never aborts).
func (m *Migrator) prepare(ctx context.Context, force bool) ([]Migration, map[string]Record, error) {
	if err := m.store.ensureTable(ctx); err != nil {
		return nil, nil, err
	}
	migrations, err := Load(m.dir)
	if err != nil {
		return nil, nil, err
	}
	records, err := m.store.list(ctx)
	if err != nil {
		return nil, nil, err
	}

	var mismatched []string
	for _, mig := range migrations {
		rec, ok := records[mig.Name]
		if !ok || rec.Status != "up" {
			continue
		}
		sum, err := checksum(mig.UpPath)
		if err != nil {
			return nil, nil, err
		}
		if sum != rec.Checksum {
			mismatched = append(mismatched, mig.Name)
		}
	}
	if len(mismatched) > 0 {
		msg := fmt.Sprintf("migrate: applied migration(s) changed on disk since they ran: %v", mismatched)
		if !force {
			return nil, nil, fmt.Errorf("%s (use --force to proceed anyway)", msg)
		}
		fmt.Fprintf(os.Stderr, "warning: %s — proceeding because --force was given\n", msg)
	}

	return migrations, records, nil
}

// withLock runs fn on a single dedicated connection, holding the dialect's advisory lock for
// fn's entire duration.
func (m *Migrator) withLock(ctx context.Context, fn func(conn *sql.Conn) error) error {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("migrate: acquiring a connection: %w", err)
	}
	defer conn.Close()

	if err := m.dialect.Lock(ctx, conn, lockKey); err != nil {
		return fmt.Errorf("migrate: acquiring lock: %w", err)
	}
	defer m.dialect.Unlock(ctx, conn, lockKey)

	return fn(conn)
}

// runFile executes path's contents. On a transactional dialect, the whole file runs in one
// transaction that is rolled back on any error. On a non-transactional dialect (MySQL: DDL
// auto-commits), it runs directly against conn and a failure partway through is reported as
// such, since it cannot be undone.
func (m *Migrator) runFile(ctx context.Context, conn *sql.Conn, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("migrate: reading %s: %w", path, err)
	}

	if !m.dialect.SupportsTransactionalDDL() {
		if _, err := conn.ExecContext(ctx, string(contents)); err != nil {
			return fmt.Errorf("migrate: running %s (dialect does not support transactional DDL — statements before the failure may already be applied): %w", path, err)
		}
		return nil
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("migrate: beginning transaction for %s: %w", path, err)
	}
	if _, err := tx.ExecContext(ctx, string(contents)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("migrate: running %s: %w", path, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migrate: committing %s: %w", path, err)
	}
	return nil
}

// Up applies the next pending migration (all == false) or every pending migration (all ==
// true), ascending by version, all sharing one new batch number. dryRun reports what would run
// without executing anything or touching the tracking table.
func (m *Migrator) Up(ctx context.Context, all bool, dryRun bool) (*Result, error) {
	result := &Result{}
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, false)
		if err != nil {
			return err
		}

		var pending []Migration
		for _, mig := range migrations {
			rec, ok := records[mig.Name]
			if !ok || rec.Status != "up" {
				pending = append(pending, mig)
			}
		}

		if len(pending) == 0 {
			return nil
		}
		if !all {
			pending = pending[:1]
		}

		if dryRun {
			for _, mig := range pending {
				result.Applied = append(result.Applied, mig.Name)
			}
			return nil
		}

		batch, err := m.store.latestBatch(ctx)
		if err != nil {
			return err
		}
		batch++

		for _, mig := range pending {
			if err := m.runFile(ctx, conn, mig.UpPath); err != nil {
				return err
			}
			sum, err := checksum(mig.UpPath)
			if err != nil {
				return err
			}
			runCount := 1
			if rec, ok := records[mig.Name]; ok {
				runCount = rec.RunCount + 1
			}
			if err := m.store.markUpWithRunCount(ctx, mig.Name, batch, sum, runCount); err != nil {
				return err
			}
			result.Applied = append(result.Applied, mig.Name)
		}
		return nil
	})
	return result, err
}

// Down rolls back the single most-recently-applied migration (all == false), or every
// migration in the current latest batch (all == true), descending by version.
func (m *Migrator) Down(ctx context.Context, all bool, dryRun bool) (*Result, error) {
	result := &Result{}
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, false)
		if err != nil {
			return err
		}
		byName := map[string]Migration{}
		for _, mig := range migrations {
			byName[mig.Name] = mig
		}

		var upNames []string
		for name, rec := range records {
			if rec.Status == "up" {
				upNames = append(upNames, name)
			}
		}
		if len(upNames) == 0 {
			return nil
		}
		sort.Sort(sort.Reverse(sort.StringSlice(upNames)))

		var toRollBack []string
		if all {
			latestBatch := records[upNames[0]].Batch
			for _, name := range upNames {
				if records[name].Batch == latestBatch {
					toRollBack = append(toRollBack, name)
				}
			}
			sort.Sort(sort.Reverse(sort.StringSlice(toRollBack)))
		} else {
			toRollBack = upNames[:1]
		}

		if dryRun {
			result.Applied = toRollBack
			return nil
		}

		for _, name := range toRollBack {
			mig, ok := byName[name]
			if !ok {
				return fmt.Errorf("migrate: %s is recorded as applied but its .sql files are missing from %s", name, m.dir)
			}
			if err := m.runFile(ctx, conn, mig.DownPath); err != nil {
				return err
			}
			if err := m.store.markDown(ctx, name); err != nil {
				return err
			}
			result.Applied = append(result.Applied, name)
		}
		return nil
	})
	return result, err
}

// Status reports every migration's on-disk file joined with its tracking-table row, if any.
func (m *Migrator) Status(ctx context.Context) ([]StatusRow, error) {
	var rows []StatusRow
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, true) // force=true: --status must never abort
		if err != nil {
			return err
		}
		for _, mig := range migrations {
			row := StatusRow{Migration: mig, ChecksumMatches: true}
			if rec, ok := records[mig.Name]; ok {
				r := rec
				row.Record = &r
				if rec.Status == "up" {
					sum, err := checksum(mig.UpPath)
					if err != nil {
						return err
					}
					row.ChecksumMatches = sum == rec.Checksum
				}
			}
			rows = append(rows, row)
		}
		return nil
	})
	return rows, err
}

// Baseline marks every migration up to and including version as already applied, without
// running any SQL. It refuses to run on a non-empty tracking table unless force is true.
func (m *Migrator) Baseline(ctx context.Context, version string, force bool) error {
	return m.withLock(ctx, func(conn *sql.Conn) error {
		if err := m.store.ensureTable(ctx); err != nil {
			return err
		}
		empty, err := m.store.isEmpty(ctx)
		if err != nil {
			return err
		}
		if !empty && !force {
			return fmt.Errorf("migrate: --baseline refuses to run on a non-empty %s table (use --force to overwrite)", m.table)
		}

		migrations, err := Load(m.dir)
		if err != nil {
			return err
		}

		found := false
		for _, mig := range migrations {
			if mig.Version == version {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("migrate: --baseline=%s does not match any migration's version", version)
		}

		for _, mig := range migrations {
			if mig.Version > version {
				continue
			}
			sum, err := checksum(mig.UpPath)
			if err != nil {
				return err
			}
			if err := m.store.baselineMark(ctx, mig.Name, sum); err != nil {
				return err
			}
		}
		return nil
	})
}
