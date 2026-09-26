package migrate

import (
	"context"
	"database/sql"
)

// fakeDialect is shared by this package's tests (store_test.go, migrate_test.go, cli_test.go).
// It delegates Open to a func field so each test supplies its own sqlmock-backed *sql.DB.
type fakeDialect struct {
	name          string
	openFn        func() (*sql.DB, error)
	transactional bool
	lockErr       error
}

func (f *fakeDialect) Name() string { return f.name }
func (f *fakeDialect) Open() (*sql.DB, error) {
	if f.openFn == nil {
		return nil, sql.ErrConnDone
	}
	return f.openFn()
}
func (f *fakeDialect) CreateMigrationsTableSQL(table string) string {
	return "CREATE TABLE IF NOT EXISTS " + table + " (id INTEGER PRIMARY KEY)"
}
func (f *fakeDialect) UpsertMigrationSQL(table string) string {
	return "INSERT INTO " + table + " (migration, batch, status, run_count, checksum, migrated_at) VALUES (?, ?, 'up', ?, ?, CURRENT_TIMESTAMP)"
}
func (f *fakeDialect) Placeholder(argPos int) string { return "?" }
func (f *fakeDialect) Lock(ctx context.Context, conn *sql.Conn, key string) error {
	return f.lockErr
}
func (f *fakeDialect) Unlock(ctx context.Context, conn *sql.Conn, key string) error { return nil }
func (f *fakeDialect) SupportsTransactionalDDL() bool                               { return f.transactional }
