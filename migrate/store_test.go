package migrate

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStore_MarkUp_FirstApplication_Inserts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO migrations`).
		WithArgs("000001_create_users", 1, 1, "abc123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	if err := s.markUp(context.Background(), "000001_create_users", 1, "abc123"); err != nil {
		t.Fatalf("markUp returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestStore_MarkUp_Reapplication_WritesGivenRunCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO migrations`).
		WithArgs("000001_create_users", 2, 2, "def456").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	if err := s.markUpWithRunCount(context.Background(), "000001_create_users", 2, "def456", 2); err != nil {
		t.Fatalf("markUpWithRunCount returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestStore_MarkDown_UpdatesStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE migrations SET status`).
		WithArgs("000001_create_users").
		WillReturnResult(sqlmock.NewResult(0, 1))

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	if err := s.markDown(context.Background(), "000001_create_users"); err != nil {
		t.Fatalf("markDown returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestStore_LatestBatch_NoRows_ReturnsZero(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(0)
	mock.ExpectQuery(`SELECT COALESCE\(MAX\(batch\), 0\)`).WillReturnRows(rows)

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	batch, err := s.latestBatch(context.Background())
	if err != nil {
		t.Fatalf("latestBatch returned error: %v", err)
	}
	if batch != 0 {
		t.Fatalf("got batch %d, want 0", batch)
	}
}

func TestStore_IsEmpty_ReflectsRowCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM migrations`).WillReturnRows(sqlmock.NewRows([]string{"c"}).AddRow(0))

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	empty, err := s.isEmpty(context.Background())
	if err != nil {
		t.Fatalf("isEmpty returned error: %v", err)
	}
	if !empty {
		t.Fatal("got empty=false, want true for a zero-row table")
	}
}

func TestStore_EnsureTable_RunsDialectSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS migrations`).WillReturnResult(sqlmock.NewResult(0, 0))

	s := newStore(db, &fakeDialect{name: "fake"}, "migrations")
	if err := s.ensureTable(context.Background()); err != nil {
		t.Fatalf("ensureTable returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
