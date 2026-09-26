package migrate

import (
	"context"
	"database/sql"
	"sort"

	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
)

// lockKey identifies this package's advisory lock to the database server. It is a fixed,
// package-level value (not derived from the tracking table name) because a single database
// should only ever run one db:migrate at a time regardless of which table it tracks into.
const lockKey = "gfly_db_migrate"

// Result reports what a single Up or Down call did, in the order it happened. An empty Applied
// means there was nothing pending (Up) or nothing applied to roll back (Down) — not an error.
type Result struct {
	Applied []string
}

// Migrator runs the *.sql files in a directory against a database, tracking each one's state in
// a table it manages itself. Construct one with NewMigrator; a Migrator is cheap to create and
// holds no long-lived connection until a method is called.
//
// A Migrator is not safe for concurrent use by multiple goroutines in the same process — call
// its methods sequentially. Concurrent *processes* (e.g. two deploys racing) are handled by the
// Dialect's advisory lock, which every method holds for its full duration.
type Migrator struct {
	db      *sql.DB
	dialect Dialect
	dir     string
	table   string
	store   *store
}

// NewMigrator constructs a Migrator that runs the migration files in dir against db, using
// dialect for every dialect-specific SQL fragment and table as the tracking table's name.
//
// table is normally DefaultTable ("migrations"); RunCLI always passes DefaultTable, since the
// spec fixes the table name rather than making it configurable (a health check and db:migrate
// must agree on one name — see the migrate package's README).
func NewMigrator(db *sql.DB, dialect Dialect, dir string, table string) *Migrator {
	return &Migrator{
		db:      db,
		dialect: dialect,
		dir:     dir,
		table:   table,
		store:   newStore(db, dialect, table),
	}
}

// Up applies pending migrations, ascending by version. With all == false (the default), it
// applies exactly one — the earliest pending migration. With all == true, it applies every
// pending migration, all sharing one freshly incremented batch number. With dryRun == true, it
// reports which migration(s) it would apply without executing any SQL or writing to the
// tracking table.
//
// Returns ErrChecksumMismatch if a migration currently marked applied no longer matches its
// on-disk checksum (see prepare). If a file fails partway through a multi-file --all run, Up
// stops immediately: earlier files in the run stay applied and recorded, the failing file is
// not recorded as applied, and no later file is attempted.
func (m *Migrator) Up(ctx context.Context, all bool, dryRun bool) (*Result, error) {
	result := &Result{Applied: []string{}}
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, false)
		if err != nil {
			return err
		}

		pending := selectPending(migrations, records)
		if len(pending) == 0 {
			return nil
		}
		if !all {
			pending = pending[:1]
		}
		if dryRun {
			result.Applied = migrationNames(pending)
			return nil
		}

		batch, err := m.store.latestBatch(ctx)
		if err != nil {
			return err
		}
		batch++

		for _, mig := range pending {
			if err := m.applyOne(ctx, conn, mig, batch, records[mig.Name].RunCount); err != nil {
				return err
			}
			result.Applied = append(result.Applied, mig.Name)
		}
		return nil
	})
	return result, err
}

// applyOne runs one migration's .up.sql file and, on success, records it as applied with
// runCount+1 (runCount is the migration's prior Record.RunCount, or 0 if it has never run).
func (m *Migrator) applyOne(ctx context.Context, conn *sql.Conn, mig Migration, batch int, priorRunCount int) error {
	if err := m.runFile(ctx, conn, mig.UpPath); err != nil {
		return err
	}
	sum, err := checksum(mig.UpPath)
	if err != nil {
		return err
	}
	return m.store.markUp(ctx, mig.Name, batch, sum, priorRunCount+1)
}

// selectPending returns every migration not currently marked StatusUp, in migrations' order
// (Load already sorts that ascending by version).
func selectPending(migrations []Migration, records map[string]Record) []Migration {
	pending := make([]Migration, 0, len(migrations))
	for _, mig := range migrations {
		if rec, ok := records[mig.Name]; !ok || rec.Status != StatusUp {
			pending = append(pending, mig)
		}
	}
	return pending
}

// Down rolls back applied migrations, descending by version. With all == false (the default),
// it rolls back exactly one — the single most-recently-applied migration, regardless of which
// batch it belongs to. With all == true, it rolls back every migration in the current latest
// batch (every StatusUp row sharing the highest Record.Batch value). With dryRun == true, it
// reports which migration(s) it would roll back without executing any SQL.
//
// Returns ErrChecksumMismatch under the same condition as Up. Returns ErrFilesMissing if a
// migration recorded as applied has no matching files in the migrations directory. As with Up,
// a failure partway through a multi-file --all rollback stops the run immediately.
func (m *Migrator) Down(ctx context.Context, all bool, dryRun bool) (*Result, error) {
	result := &Result{Applied: []string{}}
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, false)
		if err != nil {
			return err
		}

		targets := selectRollbackTargets(records, all)
		if len(targets) == 0 {
			return nil
		}
		if dryRun {
			result.Applied = targets
			return nil
		}

		byName := migrationsByName(migrations)
		for _, name := range targets {
			mig, ok := byName[name]
			if !ok {
				return errors.New("%w: %s (looked in %s)", ErrFilesMissing, name, m.dir)
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

// selectRollbackTargets returns the names of the migrations Down should roll back, descending
// by version: just the single highest-version StatusUp row when all is false, or every
// StatusUp row sharing that row's batch number when all is true. Returns an empty slice (never
// nil) if nothing is currently applied.
func selectRollbackTargets(records map[string]Record, all bool) []string {
	applied := make([]string, 0, len(records))
	for name, rec := range records {
		if rec.Status == StatusUp {
			applied = append(applied, name)
		}
	}
	if len(applied) == 0 {
		return applied
	}
	sort.Sort(sort.Reverse(sort.StringSlice(applied)))

	if !all {
		return applied[:1]
	}

	latestBatch := records[applied[0]].Batch
	targets := make([]string, 0, len(applied))
	for _, name := range applied {
		if records[name].Batch == latestBatch {
			targets = append(targets, name)
		}
	}
	return targets
}

// Status reports every migration file's state: joined with its tracking-table row if it has
// one, or reported as pending if it doesn't. Unlike Up and Down, Status never fails on
// ErrChecksumMismatch — it reports the mismatch per row (StatusRow.ChecksumMatches) instead of
// refusing to run, since inspecting state should always be safe.
func (m *Migrator) Status(ctx context.Context) ([]StatusRow, error) {
	rows := []StatusRow{}
	err := m.withLock(ctx, func(conn *sql.Conn) error {
		migrations, records, err := m.prepare(ctx, true) // force=true: Status must never abort
		if err != nil {
			return err
		}
		for _, mig := range migrations {
			rows = append(rows, statusRowFor(mig, records))
		}
		return nil
	})
	return rows, err
}

// statusRowFor builds mig's StatusRow from records, recomputing the on-disk checksum only when
// mig is currently applied (a pending migration has nothing recorded to compare against).
func statusRowFor(mig Migration, records map[string]Record) StatusRow {
	row := StatusRow{Migration: mig, ChecksumMatches: true}
	rec, ok := records[mig.Name]
	if !ok {
		return row
	}
	row.Record = &rec
	if rec.Status == StatusUp {
		if sum, err := checksum(mig.UpPath); err == nil {
			row.ChecksumMatches = sum == rec.Checksum
		}
	}
	return row
}

// Baseline marks every migration up to and including version as already applied — batch 0,
// run count 1, real checksum — without running any SQL. Use it once, against a database that
// already has this schema from before it was managed by this tool.
//
// Returns ErrNonEmptyTable if the tracking table already has any rows and force is false.
// Returns ErrUnknownVersion if version does not match any migration Load discovers. With
// force == true, Baseline proceeds even over an existing baseline, overwriting the affected
// rows.
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
			return errors.New("%w: %s", ErrNonEmptyTable, m.table)
		}

		migrations, err := Load(m.dir)
		if err != nil {
			return err
		}
		if !hasVersion(migrations, version) {
			return errors.New("%w: %s", ErrUnknownVersion, version)
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

// prepare ensures the tracking table exists, loads the migration files from disk, loads the
// current tracking rows, and verifies every StatusUp migration's on-disk checksum still matches
// what was recorded when it last ran. Every exported method calls this first; Status is the one
// caller that passes force == true, since it must report a mismatch rather than abort on it.
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

	mismatched, err := checksumMismatches(migrations, records)
	if err != nil {
		return nil, nil, err
	}
	if len(mismatched) > 0 {
		if !force {
			return nil, nil, errors.New("%w: %v", ErrChecksumMismatch, mismatched)
		}
		log.Warnf("%v: %v — proceeding because --force was given", ErrChecksumMismatch, mismatched)
	}

	return migrations, records, nil
}

// checksumMismatches returns the names of every StatusUp migration whose current on-disk
// checksum no longer matches the one recorded in its Record.
func checksumMismatches(migrations []Migration, records map[string]Record) ([]string, error) {
	var mismatched []string
	for _, mig := range migrations {
		rec, ok := records[mig.Name]
		if !ok || rec.Status != StatusUp {
			continue
		}
		sum, err := checksum(mig.UpPath)
		if err != nil {
			return nil, err
		}
		if sum != rec.Checksum {
			mismatched = append(mismatched, mig.Name)
		}
	}
	return mismatched, nil
}

// withLock runs fn on a single dedicated connection, holding the dialect's advisory lock for
// fn's entire duration, so a concurrent db:migrate process (on any connection) blocks instead of
// racing this one.
func (m *Migrator) withLock(ctx context.Context, fn func(conn *sql.Conn) error) error {
	conn, err := m.db.Conn(ctx)
	if err != nil {
		return errors.New("acquiring a connection: %w", err)
	}
	defer conn.Close()

	if err := m.dialect.Lock(ctx, conn, lockKey); err != nil {
		return errors.New("acquiring lock: %w", err)
	}
	defer m.dialect.Unlock(ctx, conn, lockKey)

	return fn(conn)
}

// runFile executes path's contents against conn. On a dialect that supports transactional DDL
// (Postgres), the whole file runs in one transaction that is rolled back on any error, so a
// failing file leaves no partial effect. On a dialect that does not (MySQL: DDL statements
// auto-commit individually), it runs directly against conn, and a failure partway through means
// the statements before it have already taken effect — the returned error says so explicitly,
// since that can't be undone by this package.
//
// It reads path with gflydev/core/utils.ReadFileAsString rather than a raw os.ReadFile — see
// the same note on source.go's checksum for why that's the right amount of abstraction here
// (small SQL files, no need for the full gflydev/storage disk abstraction).
func (m *Migrator) runFile(ctx context.Context, conn *sql.Conn, path string) error {
	contents, err := utils.ReadFileAsString(path)
	if err != nil {
		return errors.New("reading %s: %w", path, err)
	}

	if !m.dialect.SupportsTransactionalDDL() {
		if _, err := conn.ExecContext(ctx, contents); err != nil {
			return errors.New("running %s (dialect has no transactional DDL — statements before the failure may already be applied): %w", path, err)
		}
		return nil
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("beginning transaction for %s: %w", path, err)
	}
	if _, err := tx.ExecContext(ctx, contents); err != nil {
		_ = tx.Rollback()
		return errors.New("running %s: %w", path, err)
	}
	if err := tx.Commit(); err != nil {
		return errors.New("committing %s: %w", path, err)
	}
	return nil
}

// migrationsByName indexes migrations by their Name for lookup by the names selectRollbackTargets
// returns.
func migrationsByName(migrations []Migration) map[string]Migration {
	byName := make(map[string]Migration, len(migrations))
	for _, mig := range migrations {
		byName[mig.Name] = mig
	}
	return byName
}

// migrationNames extracts each Migration's Name, in order.
func migrationNames(migrations []Migration) []string {
	names := make([]string, len(migrations))
	for i, mig := range migrations {
		names[i] = mig.Name
	}
	return names
}

// hasVersion reports whether any migration in migrations has the given version.
func hasVersion(migrations []Migration, version string) bool {
	for _, mig := range migrations {
		if mig.Version == version {
			return true
		}
	}
	return false
}
