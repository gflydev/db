package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gflydev/core/utils"
)

func TestLoad_ValidDirectory_ReturnsSortedMigrations(t *testing.T) {
	migrations, err := Load(filepath.Join("testdata", "valid"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(migrations) != 2 {
		t.Fatalf("got %d migrations, want 2", len(migrations))
	}
	if migrations[0].Name != "20260101_000001_create_users" || migrations[1].Name != "20260101_000002_add_email_index" {
		t.Fatalf("got names %q, %q", migrations[0].Name, migrations[1].Name)
	}
	if migrations[0].Version != "20260101_000001" {
		t.Fatalf("got version %q, want %q", migrations[0].Version, "20260101_000001")
	}
}

func TestLoad_OrphanUpFile_ReturnsError(t *testing.T) {
	_, err := Load(filepath.Join("testdata", "orphan_up"))
	if err == nil {
		t.Fatal("expected an error for an unpaired .up.sql, got nil")
	}
}

func TestLoad_BadFilename_ReturnsError(t *testing.T) {
	_, err := Load(filepath.Join("testdata", "bad_name"))
	if err == nil {
		t.Fatal("expected an error for a filename with no up/down suffix, got nil")
	}
}

func TestLoad_OldSequentialNumberStyle_IsRejected(t *testing.T) {
	// "000001_legacy_style.up.sql" was this package's naming convention before it switched to
	// a UTC timestamp prefix. Load must reject it outright, not accept it as a differently-
	// shaped but still-valid version — a directory mixing both styles is exactly the ambiguous
	// state the strict single pattern is meant to prevent.
	_, err := Load(filepath.Join("testdata", "old_style"))
	if err == nil {
		t.Fatal("expected an error for the old NNNNNN_name.up.sql style, got nil")
	}
}

func TestChecksum_SameContent_SameChecksum(t *testing.T) {
	a, err := checksum(filepath.Join("testdata", "valid", "20260101_000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	b, err := checksum(filepath.Join("testdata", "valid", "20260101_000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	if a != b || a == "" {
		t.Fatalf("got checksums %q and %q, want equal and non-empty", a, b)
	}
}

func TestChecksum_DifferentContent_DifferentChecksum(t *testing.T) {
	a, err := checksum(filepath.Join("testdata", "valid", "20260101_000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	b, err := checksum(filepath.Join("testdata", "valid", "20260101_000002_add_email_index.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	if a == b {
		t.Fatalf("expected different checksums for different content, got %q for both", a)
	}
}

func TestNew_ValidDescription_CreatesLoadableFiles(t *testing.T) {
	dir := t.TempDir()

	mig, err := New(dir, "create_widgets_table")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if !strings.HasSuffix(mig.Name, "_create_widgets_table") {
		t.Fatalf("got name %q, want it to end with _create_widgets_table", mig.Name)
	}
	if len(mig.Version) != versionLength {
		t.Fatalf("got version %q of length %d, want length %d", mig.Version, len(mig.Version), versionLength)
	}
	if !utils.FileExists(mig.UpPath) || !utils.FileExists(mig.DownPath) {
		t.Fatalf("expected both %s and %s to exist", mig.UpPath, mig.DownPath)
	}

	// The file New just created must itself satisfy Load's own pattern — proving the two
	// stay in sync instead of drifting apart.
	migrations, err := Load(dir)
	if err != nil {
		t.Fatalf("Load of New's own output returned error: %v", err)
	}
	if len(migrations) != 1 || migrations[0].Name != mig.Name {
		t.Fatalf("got %+v, want exactly the migration New created", migrations)
	}
}

func TestNew_InvalidDescription_ReturnsErrorWithoutCreatingFiles(t *testing.T) {
	dir := t.TempDir()

	if _, err := New(dir, "Not Snake Case!"); err == nil {
		t.Fatal("expected an error for a non-snake_case description, got nil")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no files created after a rejected description, got %v", entries)
	}
}

func TestNew_FileAlreadyExists_ReturnsError(t *testing.T) {
	dir := t.TempDir()

	if _, err := New(dir, "create_widgets_table"); err != nil {
		t.Fatalf("first New returned error: %v", err)
	}
	// Calling New again for the same description within the same second collides on the
	// timestamp; simulate that by pre-creating the file a second call would target.
	collidingUp := filepath.Join(dir, time.Now().UTC().Format(timestampLayout)+"_create_widgets_table.up.sql")
	if err := os.WriteFile(collidingUp, []byte("existing"), 0o644); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	if _, err := New(dir, "create_widgets_table"); err == nil {
		t.Fatal("expected an error when the target file already exists, got nil")
	}
}
