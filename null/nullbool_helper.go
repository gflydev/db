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
// It accepts both bool and *bool values.
func Bool(val any) sql.NullBool {
	switch v := val.(type) {
	case bool:
		return sql.NullBool{
			Bool:  v,
			Valid: true,
		}
	case *bool:
		if v == nil {
			return sql.NullBool{
				Bool:  false,
				Valid: false,
			}
		}
		return sql.NullBool{
			Bool:  *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullBool
		return sql.NullBool{
			Bool:  false,
			Valid: false,
		}
	}
}

// BoolVal function will scan NullBool value.
func BoolVal(nullBool sql.NullBool) *bool {
	if !nullBool.Valid {
		return nil
	}
	return &nullBool.Bool
}
