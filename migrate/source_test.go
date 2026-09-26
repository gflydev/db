package migrate

import (
	"path/filepath"
	"testing"
)

func TestLoad_ValidDirectory_ReturnsSortedMigrations(t *testing.T) {
	migrations, err := Load(filepath.Join("testdata", "valid"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(migrations) != 2 {
		t.Fatalf("got %d migrations, want 2", len(migrations))
	}
	if migrations[0].Name != "000001_create_users" || migrations[1].Name != "000002_add_email_index" {
		t.Fatalf("got names %q, %q", migrations[0].Name, migrations[1].Name)
	}
	if migrations[0].Version != "000001" {
		t.Fatalf("got version %q, want %q", migrations[0].Version, "000001")
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

func TestChecksum_SameContent_SameChecksum(t *testing.T) {
	a, err := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	b, err := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	if a != b || a == "" {
		t.Fatalf("got checksums %q and %q, want equal and non-empty", a, b)
	}
}

func TestChecksum_DifferentContent_DifferentChecksum(t *testing.T) {
	a, err := checksum(filepath.Join("testdata", "valid", "000001_create_users.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	b, err := checksum(filepath.Join("testdata", "valid", "000002_add_email_index.up.sql"))
	if err != nil {
		t.Fatalf("checksum returned error: %v", err)
	}
	if a == b {
		t.Fatalf("expected different checksums for different content, got %q for both", a)
	}
}
