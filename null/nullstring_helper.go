package null

import (
	"database/sql"
	"database/sql/driver"
)

// ScanString function will scan NullString value.
func ScanString(nullString sql.NullString) driver.Value {
	if !nullString.Valid {
		return nil
	}
	return nullString.String
}

// String function will create a NullString object.
func String(val string) sql.NullString {
	return sql.NullString{
		String: val,
		Valid:  true,
	}
}

// StringVal function will scan NullString value.
func StringVal(nullString sql.NullString) *string {
	if !nullString.Valid {
		return nil
	}
	return &nullString.String
}
