package null

import (
	"database/sql"
	"database/sql/driver"
)

// ScanBool function will scan NullBool value.
func ScanBool(nullBool sql.NullBool) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Bool
}

// Bool function will create a NullBool object.
func Bool(val bool) sql.NullBool {
	return sql.NullBool{
		Bool:  val,
		Valid: true,
	}
}

// BoolVal function will scan NullBool value.
func BoolVal(nullBool sql.NullBool) *bool {
	if !nullBool.Valid {
		return nil
	}
	return &nullBool.Bool
}
