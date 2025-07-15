package null

import (
	"database/sql"
	"database/sql/driver"
)

// ScanInt64 function will scan NullInt64 value.
func ScanInt64(nullInt sql.NullInt64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int64
}

// Int64 function will create a NullInt64 object.
// It accepts both int64 and *int64 values.
func Int64(val any) sql.NullInt64 {
	switch v := val.(type) {
	case int64:
		return sql.NullInt64{
			Int64: v,
			Valid: true,
		}
	case *int64:
		if v == nil {
			return sql.NullInt64{
				Int64: 0,
				Valid: false,
			}
		}
		return sql.NullInt64{
			Int64: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt64
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
}

// ScanInt32 function will scan NullInt32 value.
func ScanInt32(nullInt sql.NullInt32) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int32
}

// Int32 function will create a NullInt32 object.
// It accepts both int32 and *int32 values.
func Int32(val any) sql.NullInt32 {
	switch v := val.(type) {
	case int32:
		return sql.NullInt32{
			Int32: v,
			Valid: true,
		}
	case *int32:
		if v == nil {
			return sql.NullInt32{
				Int32: 0,
				Valid: false,
			}
		}
		return sql.NullInt32{
			Int32: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt32
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}
}

// ScanInt16 function will scan NullInt16 value.
func ScanInt16(nullInt sql.NullInt16) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int16
}

// Int16 function will create a NullInt16 object.
// It accepts both int16 and *int16 values.
func Int16(val any) sql.NullInt16 {
	switch v := val.(type) {
	case int16:
		return sql.NullInt16{
			Int16: v,
			Valid: true,
		}
	case *int16:
		if v == nil {
			return sql.NullInt16{
				Int16: 0,
				Valid: false,
			}
		}
		return sql.NullInt16{
			Int16: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt16
		return sql.NullInt16{
			Int16: 0,
			Valid: false,
		}
	}
}

// Int64Val function will scan NullInt64 value.
func Int64Val(nullInt sql.NullInt64) *int64 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int64
}

// Int32Val function will scan NullInt32 value.
func Int32Val(nullInt sql.NullInt32) *int32 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int32
}

// Int16Val function will scan NullInt16 value.
func Int16Val(nullInt sql.NullInt16) *int16 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int16
}
