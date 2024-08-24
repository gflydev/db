package null

import (
	"database/sql"
	"database/sql/driver"
)

// ScanFloat64 function will scan NullFloat64 value.
func ScanFloat64(nullInt sql.NullFloat64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Float64
}

// Float64 function will create a NullFloat64 object.
func Float64(val float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: val,
		Valid:   true,
	}
}
