package postgres

import (
	"net/url"
	"testing"
)

func TestDialect_CreateMigrationsTableSQL_UsesSerial(t *testing.T) {
	sql := (Dialect{}).CreateMigrationsTableSQL("migrations")
	if sql == "" {
		t.Fatal("expected a non-empty CREATE TABLE statement")
	}
	if want := "SERIAL"; !contains(sql, want) {
		t.Fatalf("expected the statement to contain %q, got: %s", want, sql)
	}
	if want := "migrations"; !contains(sql, want) {
		t.Fatalf("expected the statement to name the table, got: %s", sql)
	}
}

func TestDialect_Placeholder_UsesDollarNumbering(t *testing.T) {
	d := Dialect{}
	if got := d.Placeholder(1); got != "$1" {
		t.Fatalf("got %q, want %q", got, "$1")
	}
	if got := d.Placeholder(3); got != "$3" {
		t.Fatalf("got %q, want %q", got, "$3")
	}
}

func TestDialect_UpsertMigrationSQL_UsesOnConflict(t *testing.T) {
	sql := (Dialect{}).UpsertMigrationSQL("migrations")
	if want := "ON CONFLICT"; !contains(sql, want) {
		t.Fatalf("expected %q in the upsert statement, got: %s", want, sql)
	}
}

func TestDialect_SupportsTransactionalDDL_IsTrue(t *testing.T) {
	if !(Dialect{}).SupportsTransactionalDDL() {
		t.Fatal("Postgres supports transactional DDL")
	}
}

func TestDialect_Name(t *testing.T) {
	if got := (Dialect{}).Name(); got != "postgres" {
		t.Fatalf("got %q, want %q", got, "postgres")
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}

func TestConnURL_EscapesSpecialCharsInCredentials(t *testing.T) {
	t.Setenv("DB_USERNAME", "dance@fit")
	t.Setenv("DB_PASSWORD", "pa@ss:w/o?r#d%1")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "dancefitvn")
	t.Setenv("DB_SSL_MODE", "disable")

	u, err := url.Parse(connURL())
	if err != nil {
		t.Fatalf("connURL is not a valid URL: %v", err)
	}
	if got := u.User.Username(); got != "dance@fit" {
		t.Fatalf("username = %q", got)
	}
	if got, _ := u.User.Password(); got != "pa@ss:w/o?r#d%1" {
		t.Fatalf("password = %q", got)
	}
	if u.Hostname() != "localhost" || u.Port() != "5432" {
		t.Fatalf("host = %q port = %q", u.Hostname(), u.Port())
	}
	if u.Path != "/dancefitvn" || u.Query().Get("sslmode") != "disable" {
		t.Fatalf("path = %q query = %q", u.Path, u.RawQuery)
	}
}
