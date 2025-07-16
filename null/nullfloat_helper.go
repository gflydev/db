package null

import (
	"database/sql"
	"database/sql/driver"
)

// floatType is a constraint interface that allows either float64 or *float64 types.
// It is used in generic functions to handle both direct float64 values and pointers
// to float64 values in a type-safe manner.
type floatType interface {
	float64 | *float64
}

// Float64Any converts a sql.NullFloat64 to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable float64 value
// to database driver operations.
//
// Parameters:
//   - nullInt (sql.NullFloat64): The nullable float64 value to convert.
//
// Returns:
//   - driver.Value: The float64 value if valid, or nil if the NullFloat64 is invalid/null.
//
// Example:
//
//	nullFloat := sql.NullFloat64{Float64: 3.14, Valid: true}
//	value := Float64Any(nullFloat) // Returns: 3.14
//
//	invalidFloat := sql.NullFloat64{Float64: 0, Valid: false}
//	value := Float64Any(invalidFloat) // Returns: nil
func Float64Any(nullInt sql.NullFloat64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Float64
}

// FloatNil converts a sql.NullFloat64 to a pointer to float64 (*float64).
// This function is useful when you need to work with nullable float64 values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullFloat (sql.NullFloat64): The nullable float64 value to convert.
//
// Returns:
//   - *float64: A pointer to the float64 value if valid, or nil if the NullFloat64 is invalid/null.
//
// Example:
//
//	nullFloat := sql.NullFloat64{Float64: 3.14, Valid: true}
//	ptr := FloatNil(nullFloat) // Returns: &3.14
//
//	invalidFloat := sql.NullFloat64{Float64: 0, Valid: false}
//	ptr := FloatNil(invalidFloat) // Returns: nil
func FloatNil(nullFloat sql.NullFloat64) *float64 {
	if !nullFloat.Valid {
		return nil
	}
	return &nullFloat.Float64
}

// Float64 creates a sql.NullFloat64 from type-constrained input types.
// This function provides a type-safe way to create nullable float64 values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - float64: Creates a valid NullFloat64 with the given float64 value
//   - *float64: Creates a valid NullFloat64 from pointer (nil pointer creates invalid NullFloat64)
//
// Returns:
//   - sql.NullFloat64: A NullFloat64 struct with appropriate Valid flag and Float64 value.
//
// Examples:
//
//	// From float64 value
//	nullFloat := Float64(3.14) // Returns: {Float64: 3.14, Valid: true}
//
//	// From float64 pointer
//	floatPtr := 2.71
//	nullFloat := Float64(&floatPtr) // Returns: {Float64: 2.71, Valid: true}
//
//	// From nil pointer
//	nullFloat := Float64((*float64)(nil)) // Returns: {Float64: 0, Valid: false}
func Float64[T floatType](val T) sql.NullFloat64 {
	switch v := any(val).(type) {
	case float64:
		return sql.NullFloat64{
			Float64: v,
			Valid:   true,
		}
	case *float64:
		if v == nil {
			return sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			}
		}
		return sql.NullFloat64{
			Float64: *v,
			Valid:   true,
		}
	default:
		// For any other type, return invalid NullFloat64
		return sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
}
