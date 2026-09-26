package migrate

import (
	"bytes"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunCLI_BareHelpArgument_PrintsHelpAndExitsZero(t *testing.T) {
	// "help" has no leading dash, so the flag package would otherwise silently treat it as a
	// non-flag argument and fall through to the default (apply) action — this is the exact
	// case a user typing `db:migrate help` hits, and RunCLI must special-case it before ever
	// reaching fs.Parse.
	var stdout, stderr bytes.Buffer
	code := runCLI([]string{"help"}, &fakeDialect{name: "fake"}, "testdata/valid", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("got exit code %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected nothing on stderr, got: %s", stderr.String())
	}
	for _, want := range []string{"--all", "--down", "--status", "--baseline", "--force"} {
		if !bytes.Contains(stdout.Bytes(), []byte(want)) {
			t.Fatalf("expected help text to mention %q, got: %s", want, stdout.String())
		}
	}
}

func TestRunCLI_DashDashHelpFlag_PrintsHelpAndExitsZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runCLI([]string{"--help"}, &fakeDialect{name: "fake"}, "testdata/valid", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("got exit code %d, want 0", code)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("--baseline")) {
		t.Fatalf("expected help text on stdout, got: %s", stdout.String())
	}
}

func TestRunCLI_New_CreatesFilesWithoutTouchingTheDatabase(t *testing.T) {
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// The fakeDialect has no openFn, so Open() would fail — proving --new never calls it.
	code := runCLI([]string{"--new=create_widgets_table"}, &fakeDialect{name: "fake"}, dir, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("got exit code %d, want 0; stderr: %s", code, stderr.String())
	}
	migrations, err := Load(dir)
	if err != nil {
		t.Fatalf("Load(dir) after --new returned error: %v", err)
	}
	if len(migrations) != 1 || !bytes.Contains([]byte(migrations[0].Name), []byte("create_widgets_table")) {
		t.Fatalf("got %+v, want one migration named *_create_widgets_table", migrations)
	}
}

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
	code := runCLI([]string{"--baseline=20260101_000001", "--down"}, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
	if code != 1 {
		t.Fatalf("got exit code %d, want 1", code)
	}
}

func TestRunCLI_BaselineWithAll_IsRejected(t *testing.T) {
	var stderr bytes.Buffer
	code := runCLI([]string{"--baseline=20260101_000001", "--all"}, &fakeDialect{name: "fake"}, "testdata/valid", &bytes.Buffer{}, &stderr)
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
	if !bytes.Contains(stdout.Bytes(), []byte("20260101_000001_create_users")) {
		t.Fatalf("expected status output to mention 20260101_000001_create_users, got: %s", stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("pending")) {
		t.Fatalf("expected status output to show a pending row, got: %s", stdout.String())
	}
}

func TestRunCLI_Status_OrphanedRecord_ShowsFileMissingAndNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("20260101_000099_deleted_from_disk", 3, "down", 1, "whatever", nil, nil))

	dialect := &fakeDialect{name: "fake", openFn: func() (*sql.DB, error) { return db, nil }}
	var stdout bytes.Buffer
	code := runCLI([]string{"--status"}, dialect, filepath.Join("testdata", "valid"), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("got exit code %d, want 0; stdout: %s", code, stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("20260101_000099_deleted_from_disk")) {
		t.Fatalf("expected the orphaned record's name in the table, got: %s", stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("FILE MISSING")) {
		t.Fatalf("expected a FILE MISSING marker, got: %s", stdout.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte("Restore the files")) {
		t.Fatalf("expected the explanatory note, got: %s", stdout.String())
	}
}

func TestRunCLI_Up_NothingPending_PrintsNothingToDo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	c1, _ := checksum(filepath.Join("testdata", "valid", "20260101_000001_create_users.up.sql"))
	c2, _ := checksum(filepath.Join("testdata", "valid", "20260101_000002_add_email_index.up.sql"))

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT migration, batch, status`).WillReturnRows(emptyRecordRows().
		AddRow("20260101_000001_create_users", 1, "up", 1, c1, nil, nil).
		AddRow("20260101_000002_add_email_index", 1, "up", 1, c2, nil, nil))

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
