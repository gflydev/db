// Package null provides utilities for handling nullable values in database operations.
// It includes helper functions for converting between sql.NullBool and other boolean
// representations, making it easier to work with nullable boolean values in Go applications.
package null

import (
	"database/sql"
)

// boolConstraint is a constraint interface that allows either bool or *bool types.
// It is used in generic functions to handle both direct boolean values and pointers
// to boolean values in a type-safe manner.
type boolConstraint interface {
	bool | *bool
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

// Bool creates a sql.NullBool from type-constrained input types.
// This function provides a type-safe way to create nullable boolean values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - bool: Creates a valid NullBool with the given boolean value
//   - *bool: Creates a valid NullBool from pointer (nil pointer creates invalid NullBool)
//
// Returns:
//   - sql.NullBool: A NullBool struct with appropriate Valid flag and Bool value.
//
// Examples:
//
//	// From bool value
//	nullBool := BoolGeneric(true) // Returns: {Bool: true, Valid: true}
//
//	// From bool pointer
//	boolPtr := &true
//	nullBool := BoolGeneric(boolPtr) // Returns: {Bool: true, Valid: true}
//
//	// From nil pointer
//	nullBool := BoolGeneric((*bool)(nil)) // Returns: {Bool: false, Valid: false}
func Bool[T boolConstraint](val T) sql.NullBool {
	switch v := any(val).(type) {
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
