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

// ByteVal function will scan NullByte value.
func ByteVal(nullByte sql.NullByte) *byte {
	if !nullByte.Valid {
		return nil
	}
	return &nullByte.Byte
}
