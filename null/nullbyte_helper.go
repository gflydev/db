package null

import (
	"database/sql"
	"database/sql/driver"
)

// ScanByte function will scan NullByte value.
func ScanByte(nullBool sql.NullByte) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Byte
}

// Byte function will create a NullBool object.
func Byte(val byte) sql.NullByte {
	return sql.NullByte{
		Byte:  val,
		Valid: true,
	}
}
