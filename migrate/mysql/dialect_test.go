package mysql

import "testing"

func TestDialect_CreateMigrationsTableSQL_UsesAutoIncrement(t *testing.T) {
	sql := (Dialect{}).CreateMigrationsTableSQL("migrations")
	if want := "AUTO_INCREMENT"; !contains(sql, want) {
		t.Fatalf("expected %q, got: %s", want, sql)
	}
}

func TestDialect_Placeholder_IsAlwaysQuestionMark(t *testing.T) {
	d := Dialect{}
	if got := d.Placeholder(1); got != "?" {
		t.Fatalf("got %q, want %q", got, "?")
	}
	if got := d.Placeholder(4); got != "?" {
		t.Fatalf("got %q, want %q", got, "?")
	}
}

func TestDialect_UpsertMigrationSQL_UsesOnDuplicateKey(t *testing.T) {
	sql := (Dialect{}).UpsertMigrationSQL("migrations")
	if want := "ON DUPLICATE KEY UPDATE"; !contains(sql, want) {
		t.Fatalf("expected %q, got: %s", want, sql)
	}
}

func TestDialect_SupportsTransactionalDDL_IsFalse(t *testing.T) {
	if (Dialect{}).SupportsTransactionalDDL() {
		t.Fatal("MySQL does not support transactional DDL")
	}
}

func TestDialect_Name(t *testing.T) {
	if got := (Dialect{}).Name(); got != "mysql" {
		t.Fatalf("got %q, want %q", got, "mysql")
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
