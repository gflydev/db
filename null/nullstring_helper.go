package null

import (
	"database/sql"
	"database/sql/driver"
)

// stringConstraint is a constraint interface that allows either string or *string types.
// It is used in generic functions to handle both direct string values and pointers
// to string values in a type-safe manner.
type stringConstraint interface {
	string | *string
}

// StringAny converts a sql.NullString to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable string value
// to database driver operations.
//
// Parameters:
//   - nullString (sql.NullString): The nullable string value to convert.
//
// Returns:
//   - driver.Value: The string value if valid, or nil if the NullString is invalid/null.
//
// Example:
//
//	nullString := sql.NullString{String: "hello", Valid: true}
//	value := StringAny(nullString) // Returns: "hello"
//
//	invalidString := sql.NullString{String: "", Valid: false}
//	value := StringAny(invalidString) // Returns: nil
func StringAny(nullString sql.NullString) driver.Value {
	if !nullString.Valid {
		return nil
	}
	return nullString.String
}

// StringNil converts a sql.NullString to a pointer to string (*string).
// This function is useful when you need to work with nullable string values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullString (sql.NullString): The nullable string value to convert.
//
// Returns:
//   - *string: A pointer to the string value if valid, or nil if the NullString is invalid/null.
//
// Example:
//
//	nullString := sql.NullString{String: "hello", Valid: true}
//	ptr := StringNil(nullString) // Returns: &"hello"
//
//	invalidString := sql.NullString{String: "", Valid: false}
//	ptr := StringNil(invalidString) // Returns: nil
func StringNil(nullString sql.NullString) *string {
	if !nullString.Valid {
		return nil
	}
	return &nullString.String
}

// StringVal returns the string value of a sql.NullString.
// If the NullString is valid, it returns the string value.
// If the NullString is invalid, it returns an empty string.
//
// Parameters:
//   - nullString (sql.NullString): The nullable string value to convert.
//
// Returns:
//   - string: The string value if valid, or an empty string if invalid.
//
// Example:
//
//	nullString := sql.NullString{String: "hello", Valid: true}
//	result := StringVal(nullString) // Returns: "hello"
//
//	invalidString := sql.NullString{String: "", Valid: false}
//	result := StringVal(invalidString) // Returns: ""
func StringVal(nullString sql.NullString) string {
	if !nullString.Valid {
		return ""
	}

	return nullString.String
}

// String creates a sql.NullString from type-constrained input types.
// This function provides a type-safe way to create nullable string values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - string: Creates a valid NullString with the given string value
//   - *string: Creates a valid NullString from pointer (nil pointer creates invalid NullString)
//
// Returns:
//   - sql.NullString: A NullString struct with appropriate Valid flag and String value.
//
// Examples:
//
//	// From string value
//	nullString := String("hello") // Returns: {String: "hello", Valid: true}
//
//	// From string pointer
//	strPtr := &"world"
//	nullString := String(strPtr) // Returns: {String: "world", Valid: true}
//
//	// From nil pointer
//	nullString := String((*string)(nil)) // Returns: {String: "", Valid: false}
func String[T stringConstraint](val T) sql.NullString {
	switch v := any(val).(type) {
	case string:
		return sql.NullString{
			String: v,
			Valid:  true,
		}
	case *string:
		if v == nil {
			return sql.NullString{
				String: "",
				Valid:  false,
			}
		}
		return sql.NullString{
			String: *v,
			Valid:  true,
		}
	default:
		// For any other type, return invalid NullString
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
}
