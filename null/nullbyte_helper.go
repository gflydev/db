package null

import (
	"database/sql"
	"database/sql/driver"
)

// ByteAny converts a sql.NullByte to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable byte value
// to database driver operations.
//
// Parameters:
//   - nullBool (sql.NullByte): The nullable byte value to convert.
//
// Returns:
//   - driver.Value: The byte value if valid, or nil if the NullByte is invalid/null.
//
// Example:
//
//	nullByte := sql.NullByte{Byte: 65, Valid: true}
//	value := ByteAny(nullByte) // Returns: 65
//
//	invalidByte := sql.NullByte{Byte: 0, Valid: false}
//	value := ByteAny(invalidByte) // Returns: nil
func ByteAny(nullBool sql.NullByte) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Byte
}

// ByteNil converts a sql.NullByte to a pointer to byte (*byte).
// This function is useful when you need to work with nullable byte values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullByte (sql.NullByte): The nullable byte value to convert.
//
// Returns:
//   - *byte: A pointer to the byte value if valid, or nil if the NullByte is invalid/null.
//
// Example:
//
//	nullByte := sql.NullByte{Byte: 65, Valid: true}
//	ptr := ByteNil(nullByte) // Returns: &65
//
//	invalidByte := sql.NullByte{Byte: 0, Valid: false}
//	ptr := ByteNil(invalidByte) // Returns: nil
func ByteNil(nullByte sql.NullByte) *byte {
	if !nullByte.Valid {
		return nil
	}
	return &nullByte.Byte
}

// Byte creates a sql.NullByte from various input types.
// This function provides a convenient way to create nullable byte values
// for database operations, handling both direct values and pointers.
//
// Parameters:
//   - val (any): The input value to convert. Supported types:
//   - byte: Creates a valid NullByte with the given byte value
//   - *byte: Creates a valid NullByte from pointer (nil pointer creates invalid NullByte)
//   - any other type: Creates an invalid NullByte with zero value
//
// Returns:
//   - sql.NullByte: A NullByte struct with appropriate Valid flag and Byte value.
//
// Examples:
//
//	// From byte value
//	nullByte := Byte(byte(65)) // Returns: {Byte: 65, Valid: true}
//
//	// From byte pointer
//	bytePtr := byte(97)
//	nullByte := Byte(&bytePtr) // Returns: {Byte: 97, Valid: true}
//
//	// From nil pointer
//	nullByte := Byte((*byte)(nil)) // Returns: {Byte: 0, Valid: false}
//
//	// From unsupported type
//	nullByte := Byte("invalid") // Returns: {Byte: 0, Valid: false}
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
