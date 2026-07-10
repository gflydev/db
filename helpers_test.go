package db

import (
	"testing"

	"github.com/gflydev/core"
)

// TestProcessOrderByAllowList verifies ordering against a []string allow-list,
// including the descending "-" prefix and rejection of unknown keys.
func TestProcessOrderByAllowList(t *testing.T) {
	t.Run("ascending accepted key", func(t *testing.T) {
		db := Instance()
		ProcessOrderBy(db, "title", []string{"created_at", "title"})
		items := db.orderByStatement.Items
		if len(items) != 1 || items[0].Field != "title" || items[0].Direction != Asc {
			t.Errorf("got %+v, want title ASC", items)
		}
	})

	t.Run("descending prefix", func(t *testing.T) {
		db := Instance()
		ProcessOrderBy(db, "-created_at", []string{"created_at", "title"})
		items := db.orderByStatement.Items
		if len(items) != 1 || items[0].Field != "created_at" || items[0].Direction != Desc {
			t.Errorf("got %+v, want created_at DESC", items)
		}
	})

	t.Run("unknown key with default", func(t *testing.T) {
		db := Instance()
		ProcessOrderBy(db, "unknown", []string{"created_at"}, "created_at", Desc)
		items := db.orderByStatement.Items
		if len(items) != 1 || items[0].Field != "created_at" || items[0].Direction != Desc {
			t.Errorf("got %+v, want fallback created_at DESC", items)
		}
	})

	t.Run("unknown key without default adds nothing", func(t *testing.T) {
		db := Instance()
		ProcessOrderBy(db, "unknown", []string{"created_at"})
		if len(db.orderByStatement.Items) != 0 {
			t.Errorf("expected no order by, got %+v", db.orderByStatement.Items)
		}
	})
}

// TestProcessOrderByMapping verifies ordering against a core.Data column map.
func TestProcessOrderByMapping(t *testing.T) {
	db := Instance()
	ProcessOrderBy(db, "-name", core.Data{"id": "categories.id", "name": "categories.name"})
	items := db.orderByStatement.Items
	if len(items) != 1 || items[0].Field != "categories.name" || items[0].Direction != Desc {
		t.Errorf("got %+v, want categories.name DESC", items)
	}
}

// TestToQBCondition verifies conversion of a Condition, including nested groups.
func TestToQBCondition(t *testing.T) {
	cond := Condition{
		Field: "status",
		Opt:   Eq,
		Value: "active",
		AndOr: And,
		Group: []Condition{
			{Field: "age", Opt: Greater, Value: 18, AndOr: And},
		},
	}

	got := cond.ToQBCondition()
	if got.Field != "status" || got.Value != "active" {
		t.Errorf("top-level mismatch: %+v", got)
	}
	if len(got.Group) != 1 || got.Group[0].Field != "age" {
		t.Errorf("group not converted: %+v", got.Group)
	}
}
