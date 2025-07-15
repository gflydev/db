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

// Byte function will create a NullByte object.
// It accepts both byte and *byte values.
func Byte(val any) sql.NullByte {
	switch v := val.(type) {
	case byte:
		return sql.NullByte{
			Byte:  v,
			Valid: true,
		}
	case *byte:
		if v == nil {
			return sql.NullByte{
				Byte:  0,
				Valid: false,
			}
		}
		return sql.NullByte{
			Byte:  *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullByte
		return sql.NullByte{
			Byte:  0,
			Valid: false,
		}
	}
}

// ByteVal function will scan NullByte value.
func ByteVal(nullByte sql.NullByte) *byte {
	if !nullByte.Valid {
		return nil
	}
	return &nullByte.Byte
}
