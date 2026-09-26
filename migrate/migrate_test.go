package migrate

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestMigrator(db *sql.DB, dir string) *Migrator {
	return NewMigrator(db, &fakeDialect{name: "fake", transactional: true}, dir, "migrations")
}

func emptyRecordRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"migration", "batch", "status", "run_count", "checksum", "migrated_at", "rolled_back_at"})
}

func TestMigrator_Up_NoAll_AppliesOneFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows())
	mock.ExpectQuery(`SELECT COALESCE\(MAX\(batch\), 0\)`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(`CREATE TABLE users`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`INSERT INTO migrations`).WillReturnResult(sqlmock.NewResult(1, 1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Up(context.Background(), false, false)
	if err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	if len(result.Applied) != 1 || result.Applied[0] != "000001_create_users" {
		t.Fatalf("got Applied %v, want [000001_create_users]", result.Applied)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMigrator_Up_All_AppliesEveryPendingFileInOneBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows())
	mock.ExpectQuery(`SELECT COALESCE\(MAX\(batch\), 0\)`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(3))
	mock.ExpectBegin()
	mock.ExpectExec(`CREATE TABLE users`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`INSERT INTO migrations`).WithArgs("000001_create_users", 4, 1, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`CREATE INDEX idx_users_email`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`INSERT INTO migrations`).WithArgs("000002_add_email_index", 4, 1, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Up(context.Background(), true, false)
	if err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	if len(result.Applied) != 2 {
		t.Fatalf("got %d applied, want 2", len(result.Applied))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMigrator_Up_ChecksumMismatch_AbortsWithoutForce(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("000001_create_users", 1, "up", 1, "0000000000000000000000000000000000000000000000000000000000000000", nil, nil))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	_, err = m.Up(context.Background(), true, false)
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Fatalf("got error %v, want errors.Is(err, ErrChecksumMismatch)", err)
	}
}

func TestMigrator_Down_NoAll_RollsBackSingleLatest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	upChecksum, err := checksum(filepath.Join("testdata", "valid", "000002_add_email_index.up.sql"))
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}
	otherChecksum, err := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("000001_create_users", 1, "up", 1, otherChecksum, nil, nil).
		AddRow("000002_add_email_index", 1, "up", 1, upChecksum, nil, nil))
	mock.ExpectBegin()
	mock.ExpectExec(`DROP INDEX idx_users_email`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`UPDATE migrations SET status`).WillReturnResult(sqlmock.NewResult(0, 1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Down(context.Background(), false, false)
	if err != nil {
		t.Fatalf("Down returned error: %v", err)
	}
	if len(result.Applied) != 1 || result.Applied[0] != "000002_add_email_index" {
		t.Fatalf("got Applied %v, want [000002_add_email_index]", result.Applied)
	}
}

func TestMigrator_Down_All_RollsBackEntireLatestBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	c1, _ := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	c2, _ := checksum(filepath.Join("testdata", "valid", "000002_add_email_index.up.sql"))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("000001_create_users", 1, "up", 1, c1, nil, nil).
		AddRow("000002_add_email_index", 1, "up", 1, c2, nil, nil))
	mock.ExpectBegin()
	mock.ExpectExec(`DROP INDEX idx_users_email`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`UPDATE migrations SET status`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`DROP TABLE users`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectExec(`UPDATE migrations SET status`).WillReturnResult(sqlmock.NewResult(0, 1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Down(context.Background(), true, false)
	if err != nil {
		t.Fatalf("Down returned error: %v", err)
	}
	if len(result.Applied) != 2 {
		t.Fatalf("got %d rolled back, want 2 (both in batch 1)", len(result.Applied))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestMigrator_Baseline_NonEmptyTable_RefusesWithoutForce(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM migrations`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	err = m.Baseline(context.Background(), "000001", false)
	if !errors.Is(err, ErrNonEmptyTable) {
		t.Fatalf("got error %v, want errors.Is(err, ErrNonEmptyTable)", err)
	}
}

func TestMigrator_Baseline_UnknownVersion_ReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM migrations`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	err = m.Baseline(context.Background(), "999999", false)
	if !errors.Is(err, ErrUnknownVersion) {
		t.Fatalf("got error %v, want errors.Is(err, ErrUnknownVersion)", err)
	}
}

func TestMigrator_Down_RecordedMigrationHasNoFiles_ReturnsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("000099_deleted_from_disk", 1, "up", 1, "whatever-checksum-does-not-matter-when-mismatch-check-is-skipped", nil, nil))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	_, err = m.Down(context.Background(), false, false)
	if !errors.Is(err, ErrFilesMissing) {
		t.Fatalf("got error %v, want errors.Is(err, ErrFilesMissing)", err)
	}
}

func TestMigrator_Baseline_EmptyTable_MarksFilesWithoutRunningSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM migrations`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))
	// baseline covers only 000001, not 000002 — exactly one upsert expected.
	mock.ExpectExec(`INSERT INTO migrations`).WithArgs("000001_create_users", 0, 1, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	if err := m.Baseline(context.Background(), "000001", false); err != nil {
		t.Fatalf("Baseline returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations (no SQL file should have run): %v", err)
	}
}

func TestMigrator_Up_MidBatchFailure_StopsAndDoesNotMarkFailedFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows())
	mock.ExpectQuery(`SELECT COALESCE\(MAX\(batch\), 0\)`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec(`CREATE TABLE users`).WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	_, err = m.Up(context.Background(), true, false)
	if err == nil {
		t.Fatal("expected an error when the first file's exec fails, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations (second file must not have been attempted): %v", err)
	}
}

func TestMigrator_Up_DryRun_DoesNotExecuteOrWrite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows())

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Up(context.Background(), true, true)
	if err != nil {
		t.Fatalf("Up (dry-run) returned error: %v", err)
	}
	if len(result.Applied) != 2 {
		t.Fatalf("got %d planned, want 2", len(result.Applied))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations (no exec/write should have happened): %v", err)
	}
}

func TestMigrator_Up_NothingPending_ReturnsEmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	c1, _ := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	c2, _ := checksum(filepath.Join("testdata", "valid", "000002_add_email_index.up.sql"))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("000001_create_users", 1, "up", 1, c1, nil, nil).
		AddRow("000002_add_email_index", 1, "up", 1, c2, nil, nil))

	m := newTestMigrator(db, filepath.Join("testdata", "valid"))
	result, err := m.Up(context.Background(), true, false)
	if err != nil {
		t.Fatalf("Up returned error: %v", err)
	}
	if len(result.Applied) != 0 {
		t.Fatalf("got %d applied, want 0 (nothing pending)", len(result.Applied))
	}
}
