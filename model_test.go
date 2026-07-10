package db

import (
	"reflect"
	"testing"
)

// TestToSnakeCase verifies conversion from camelCase/PascalCase to snake_case.
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"pascal single word", "User", "user"},
		{"pascal two words", "UserDetail", "user_detail"},
		{"camel case", "userDetail", "user_detail"},
		{"acronym then word", "APIKey", "api_key"},
		{"trailing acronym", "UserID", "user_id"},
		{"with digits", "Address2Line", "address2_line"},
		{"already snake", "user_name", "user_name"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSnakeCase(tt.input); got != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestReadTags verifies parsing of the `model` struct tag into an attribute map.
func TestReadTags(t *testing.T) {
	t.Run("empty tag defaults to BOOLEAN", func(t *testing.T) {
		got := readTags("")
		want := map[string][]string{TYPE: {"BOOLEAN"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("readTags(\"\") = %v, want %v", got, want)
		}
	})

	t.Run("type with primary", func(t *testing.T) {
		got := readTags("type:serial,primary")
		want := map[string][]string{TYPE: {"serial", "primary"}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("table name", func(t *testing.T) {
		got := readTags("table:users")
		if v, ok := got[TABLE]; !ok || len(v) != 1 || v[0] != "users" {
			t.Errorf("got %v, want table:[users]", got)
		}
	})

	t.Run("multiple attributes with spaces", func(t *testing.T) {
		got := readTags("type: varchar(255) ; name: full_name")
		if v := got[TYPE]; len(v) != 1 || v[0] != "varchar(255)" {
			t.Errorf("type = %v, want [varchar(255)]", v)
		}
		if v := got[NAME]; len(v) != 1 || v[0] != "full_name" {
			t.Errorf("name = %v, want [full_name]", v)
		}
	})
}

// TestIsPrimaryAndIsSerial verifies attribute detection helpers.
func TestIsPrimaryAndIsSerial(t *testing.T) {
	if !isPrimary([]string{"serial", "primary"}) {
		t.Error("isPrimary should detect primary")
	}
	if isPrimary([]string{"varchar(255)"}) {
		t.Error("isPrimary should be false for non-primary")
	}
	if !isSerial([]string{"serial", "primary"}) {
		t.Error("isSerial should detect serial")
	}
	if isSerial([]string{"numeric"}) {
		t.Error("isSerial should be false for non-serial")
	}
}

// TestGetTypes verifies type-attribute formatting for schema definitions.
func TestGetTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"primary keyword expands", []string{"primary"}, "PRIMARY KEY "},
		{"uppercases plain type", []string{"varchar(255)"}, "VARCHAR(255) "},
		{"serial and primary", []string{"serial", "primary"}, "SERIAL PRIMARY KEY "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTypes(tt.input); got != tt.expected {
				t.Errorf("getTypes(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestGetCascade verifies cascade rule generation.
func TestGetCascade(t *testing.T) {
	if got := getCascade([]string{"delete"}); got != "ON DELETE CASCADE " {
		t.Errorf("getCascade(delete) = %q", got)
	}
	if got := getCascade([]string{"update"}); got != "ON UPDATE CASCADE " {
		t.Errorf("getCascade(update) = %q", got)
	}
}

// TestGetReferences verifies REFERENCES clause construction.
func TestGetReferences(t *testing.T) {
	got := getReferences("Users", "user_id")
	want := "REFERENCES users (id) "
	if got != want {
		t.Errorf("getReferences = %q, want %q", got, want)
	}
}

// sampleUser is a model used to exercise ModelData reflection.
type sampleUser struct {
	MetaData MetaData `db:"-" model:"table:users"`
	Id       int      `db:"id" model:"type:serial,primary"`
	Name     string   `db:"name" model:"type:varchar(255)"`
}

// TestModelData verifies struct-to-Table conversion including table name,
// column extraction, and primary/serial detection.
func TestModelData(t *testing.T) {
	tbl, err := ModelData(sampleUser{Id: 5, Name: "Vinh"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tbl.Name != "users" {
		t.Errorf("table name = %q, want users", tbl.Name)
	}
	// MetaData field is skipped, so only id and name are columns.
	if len(tbl.Columns) != 2 {
		t.Errorf("columns = %d, want 2", len(tbl.Columns))
	}
	if len(tbl.Primaries) != 1 || tbl.Primaries[0].Name != "id" {
		t.Errorf("primaries = %v, want [id]", tbl.Primaries)
	}
	if tbl.PrimarySerial == nil || tbl.PrimarySerial.Name != "id" {
		t.Errorf("PrimarySerial = %v, want id", tbl.PrimarySerial)
	}
	if tbl.Values["name"] != "Vinh" {
		t.Errorf("value name = %v, want Vinh", tbl.Values["name"])
	}
}

// TestModelDataInvalidInput verifies a non-struct input is rejected.
func TestModelDataInvalidInput(t *testing.T) {
	if _, err := ModelData(42); err == nil {
		t.Error("expected error for non-struct input")
	}
}
