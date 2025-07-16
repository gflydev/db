package null

import (
	"database/sql"
	"database/sql/driver"
)

// BoolAny converts a sql.NullBool to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable boolean value
// to database driver operations.
//
// Parameters:
//   - nullBool (sql.NullBool): The nullable boolean value to convert.
//
// Returns:
//   - driver.Value: The boolean value if valid, or nil if the NullBool is invalid/null.
//
// Example:
//
//	nullBool := sql.NullBool{Bool: true, Valid: true}
//	value := BoolAny(nullBool) // Returns: true
//
//	invalidBool := sql.NullBool{Bool: false, Valid: false}
//	value := BoolAny(invalidBool) // Returns: nil
func BoolAny(nullBool sql.NullBool) driver.Value {
	if !nullBool.Valid {
		return nil
	}
	return nullBool.Bool
}

// BoolNil converts a sql.NullBool to a pointer to bool (*bool).
// This function is useful when you need to work with nullable boolean values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullBool (sql.NullBool): The nullable boolean value to convert.
//
// Returns:
//   - *bool: A pointer to the boolean value if valid, or nil if the NullBool is invalid/null.
//
// Example:
//
//	nullBool := sql.NullBool{Bool: true, Valid: true}
//	ptr := BoolNil(nullBool) // Returns: &true
//
//	invalidBool := sql.NullBool{Bool: false, Valid: false}
//	ptr := BoolNil(invalidBool) // Returns: nil
func BoolNil(nullBool sql.NullBool) *bool {
	if !nullBool.Valid {
		return nil
	}
	return &nullBool.Bool
}

// Bool creates a sql.NullBool from various input types.
// This function provides a convenient way to create nullable boolean values
// for database operations, handling both direct values and pointers.
//
// Parameters:
//   - val (any): The input value to convert. Supported types:
//   - bool: Creates a valid NullBool with the given boolean value
//   - *bool: Creates a valid NullBool from pointer (nil pointer creates invalid NullBool)
//   - any other type: Creates an invalid NullBool with false value
//
// Returns:
//   - sql.NullBool: A NullBool struct with appropriate Valid flag and Bool value.
//
// Examples:
//
//	// From bool value
//	nullBool := Bool(true) // Returns: {Bool: true, Valid: true}
//
//	// From bool pointer
//	boolPtr := &true
//	nullBool := Bool(boolPtr) // Returns: {Bool: true, Valid: true}
//
//	// From nil pointer
//	nullBool := Bool((*bool)(nil)) // Returns: {Bool: false, Valid: false}
//
//	// From unsupported type
//	nullBool := Bool("invalid") // Returns: {Bool: false, Valid: false}
func Bool(val any) sql.NullBool {
	switch v := val.(type) {
	case bool:
		return sql.NullBool{
			Bool:  v,
			Valid: true,
		}
	case *bool:
		if v == nil {
			return sql.NullBool{
				Bool:  false,
				Valid: false,
			}
		}
		return sql.NullBool{
			Bool:  *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullBool
		return sql.NullBool{
			Bool:  false,
			Valid: false,
		}
	}
}
