package null

import (
	"database/sql"
	"database/sql/driver"
)

// StringAny function will scan NullString value.
func StringAny(nullString sql.NullString) driver.Value {
	if !nullString.Valid {
		return nil
	}
	return nullString.String
}

// StringNil function will scan NullString value.
func StringNil(nullString sql.NullString) *string {
	if !nullString.Valid {
		return nil
	}
	return &nullString.String
}

// String function will create a NullString object.
// It accepts both string and *string values.
func String(val any) sql.NullString {
	switch v := val.(type) {
	case string:
		return sql.NullString{
			String: v,
			Valid:  true,
		}
	case *string:
		if v == nil {
			return sql.NullString{
				String: "",
				Valid:  false,
			}
		}
		return sql.NullString{
			String: *v,
			Valid:  true,
		}
	default:
		// For any other type, return invalid NullString
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
}
