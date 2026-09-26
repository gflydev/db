//go:build integration

package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	migrate "github.com/gflydev/db/migrate"
)

// TestDialect_Integration_FullCycle exercises up -> status -> down -> up-again against a real
// Postgres. Run with DB_* env vars pointed at a disposable database:
//
//	DB_HOST=localhost DB_PORT=5432 DB_NAME=gfly_migrate_test DB_USERNAME=... DB_PASSWORD=... \
//	  go test -tags=integration ./... -v
func TestDialect_Integration_FullCycle(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "000001_create_widgets.up.sql"), "CREATE TABLE widgets (id SERIAL PRIMARY KEY);")
	writeFile(t, filepath.Join(dir, "000001_create_widgets.down.sql"), "DROP TABLE widgets;")

	db, err := (Dialect{}).Open()
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Skipf("no reachable Postgres (set DB_* env vars): %v", err)
	}
	t.Cleanup(func() {
		db.Exec("DROP TABLE IF EXISTS migrations")
		db.Exec("DROP TABLE IF EXISTS widgets")
	})

	m := migrate.NewMigrator(db, Dialect{}, dir, "migrations")
	ctx := context.Background()

	result, err := m.Up(ctx, false, false)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}
	if len(result.Applied) != 1 {
		t.Fatalf("got %d applied, want 1", len(result.Applied))
	}

	var exists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'widgets')").Scan(&exists); err != nil {
		t.Fatalf("checking widgets table: %v", err)
	}
	if !exists {
		t.Fatal("widgets table was not created")
	}

	rows, err := m.Status(ctx)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(rows) != 1 || rows[0].Record == nil || rows[0].Record.RunCount != 1 {
		t.Fatalf("got status %+v, want one row with run_count 1", rows)
	}

	if _, err := m.Down(ctx, false, false); err != nil {
		t.Fatalf("Down: %v", err)
	}
	if err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'widgets')").Scan(&exists); err != nil {
		t.Fatalf("checking widgets table after Down: %v", err)
	}
	if exists {
		t.Fatal("widgets table still exists after Down")
	}

	result, err = m.Up(ctx, false, false)
	if err != nil {
		t.Fatalf("second Up: %v", err)
	}
	if len(result.Applied) != 1 {
		t.Fatalf("got %d applied on re-run, want 1", len(result.Applied))
	}
	rows, err = m.Status(ctx)
	if err != nil {
		t.Fatalf("Status after re-run: %v", err)
	}
	if rows[0].Record.RunCount != 2 {
		t.Fatalf("got run_count %d after re-applying, want 2", rows[0].Record.RunCount)
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
