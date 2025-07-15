package null

import (
	"database/sql"
	"database/sql/driver"
)

// Float64Any function will scan NullFloat64 value.
func Float64Any(nullInt sql.NullFloat64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Float64
}

// FloatNil function will scan NullFloat64 value.
func FloatNil(nullFloat sql.NullFloat64) *float64 {
	if !nullFloat.Valid {
		return nil
	}
	return &nullFloat.Float64
}

// Float64 function will create a NullFloat64 object.
// It accepts both float64 and *float64 values.
func Float64(val any) sql.NullFloat64 {
	switch v := val.(type) {
	case float64:
		return sql.NullFloat64{
			Float64: v,
			Valid:   true,
		}
	case *float64:
		if v == nil {
			return sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			}
		}
		return sql.NullFloat64{
			Float64: *v,
			Valid:   true,
		}
	default:
		// For any other type, return invalid NullFloat64
		return sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
}
