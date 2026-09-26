package migrate

import (
	"bytes"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunCLI_UnknownFlag_ReturnsOneAndPrintsUsage(t *testing.T) {
	var stderr bytes.Buffer
	code := runCLI([]string{"--nonsense"}, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
	if code != 1 {
		t.Fatalf("got exit code %d, want 1", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected usage/error output on stderr, got none")
	}
}

func TestRunCLI_BaselineWithDown_IsRejected(t *testing.T) {
	var stderr bytes.Buffer
	code := runCLI([]string{"--baseline=000001", "--down"}, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
	if code != 1 {
		t.Fatalf("got exit code %d, want 1", code)
	}
}

func TestRunCLI_BaselineWithAll_IsRejected(t *testing.T) {
	var stderr bytes.Buffer
	code := runCLI([]string{"--baseline=000001", "--all"}, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
	if code != 1 {
		t.Fatalf("got exit code %d, want 1", code)
	}
}

func TestRunCLI_OpenFails_ReturnsOne(t *testing.T) {
	var stderr bytes.Buffer
	code := runCLI(nil, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
	if code != 1 {
		t.Fatalf("got exit code %d, want 1", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected a connection error on stderr, got none")
	}
}

func TestRunCLI_Status_PrintsOneLinePerMigration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows())

	dialect := &fakeDialect{name: "fake", openFn: func() (*sql.DB, error) { return db, nil }}
	var stdout bytes.Buffer
	code := runCLI([]string{"--status"}, dialect, filepath.Join("testdata", "valid"), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("got exit code %d, want 0; stdout: %s", code, stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("000001_create_users")) {
		t.Fatalf("expected status output to mention 000001_create_users, got: %s", stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("pending")) {
		t.Fatalf("expected status output to show a pending row, got: %s", stdout.String())
	}
}

func TestRunCLI_Up_NothingPending_PrintsNothingToDo(t *testing.T) {
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

	dialect := &fakeDialect{name: "fake", openFn: func() (*sql.DB, error) { return db, nil }}
	var stdout bytes.Buffer
	code := runCLI([]string{"--all"}, dialect, filepath.Join("testdata", "valid"), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("got exit code %d, want 0", code)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("Nothing to do")) {
		t.Fatalf("expected 'Nothing to do.', got: %s", stdout.String())
	}
}
