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
func Int64(val int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: val,
		Valid: true,
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
func Int32(val int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: val,
		Valid: true,
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
func Int16(val int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: val,
		Valid: true,
	}
}
