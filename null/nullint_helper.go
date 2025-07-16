package null

import (
	"database/sql"
	"database/sql/driver"
)

// int64Type is a constraint interface that allows either int64 or *int64 types.
// It is used in generic functions to handle both direct int64 values and pointers
// to int64 values in a type-safe manner.
type int64Type interface {
	int64 | *int64
}

// int32Type is a constraint interface that allows either int32 or *int32 types.
// It is used in generic functions to handle both direct int32 values and pointers
// to int32 values in a type-safe manner.
type int32Type interface {
	int32 | *int32
}

// int16Type is a constraint interface that allows either int16 or *int16 types.
// It is used in generic functions to handle both direct int16 values and pointers
// to int16 values in a type-safe manner.
type int16Type interface {
	int16 | *int16
}

// Int64Any converts a sql.NullInt64 to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable int64 value
// to database driver operations.
//
// Parameters:
//   - nullInt (sql.NullInt64): The nullable int64 value to convert.
//
// Returns:
//   - driver.Value: The int64 value if valid, or nil if the NullInt64 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt64{Int64: 42, Valid: true}
//	value := Int64Any(nullInt) // Returns: 42
//
//	invalidInt := sql.NullInt64{Int64: 0, Valid: false}
//	value := Int64Any(invalidInt) // Returns: nil
func Int64Any(nullInt sql.NullInt64) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int64
}

// Int64Nil converts a sql.NullInt64 to a pointer to int64 (*int64).
// This function is useful when you need to work with nullable int64 values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullInt (sql.NullInt64): The nullable int64 value to convert.
//
// Returns:
//   - *int64: A pointer to the int64 value if valid, or nil if the NullInt64 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt64{Int64: 42, Valid: true}
//	ptr := Int64Nil(nullInt) // Returns: &42
//
//	invalidInt := sql.NullInt64{Int64: 0, Valid: false}
//	ptr := Int64Nil(invalidInt) // Returns: nil
func Int64Nil(nullInt sql.NullInt64) *int64 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int64
}

// Int64 creates a sql.NullInt64 from type-constrained input types.
// This function provides a type-safe way to create nullable int64 values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - int64: Creates a valid NullInt64 with the given int64 value
//   - *int64: Creates a valid NullInt64 from pointer (nil pointer creates invalid NullInt64)
//
// Returns:
//   - sql.NullInt64: A NullInt64 struct with appropriate Valid flag and Int64 value.
//
// Examples:
//
//	// From int64 value
//	nullInt := Int64(int64(42)) // Returns: {Int64: 42, Valid: true}
//
//	// From int64 pointer
//	intPtr := int64(100)
//	nullInt := Int64(&intPtr) // Returns: {Int64: 100, Valid: true}
//
//	// From nil pointer
//	nullInt := Int64((*int64)(nil)) // Returns: {Int64: 0, Valid: false}
func Int64[T int64Type](val T) sql.NullInt64 {
	switch v := any(val).(type) {
	case int64:
		return sql.NullInt64{
			Int64: v,
			Valid: true,
		}
	case *int64:
		if v == nil {
			return sql.NullInt64{
				Int64: 0,
				Valid: false,
			}
		}
		return sql.NullInt64{
			Int64: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt64
		return sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
	}
}

// Int32Any converts a sql.NullInt32 to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable int32 value
// to database driver operations.
//
// Parameters:
//   - nullInt (sql.NullInt32): The nullable int32 value to convert.
//
// Returns:
//   - driver.Value: The int32 value if valid, or nil if the NullInt32 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt32{Int32: 42, Valid: true}
//	value := Int32Any(nullInt) // Returns: 42
//
//	invalidInt := sql.NullInt32{Int32: 0, Valid: false}
//	value := Int32Any(invalidInt) // Returns: nil
func Int32Any(nullInt sql.NullInt32) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int32
}

// Int32Nil converts a sql.NullInt32 to a pointer to int32 (*int32).
// This function is useful when you need to work with nullable int32 values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullInt (sql.NullInt32): The nullable int32 value to convert.
//
// Returns:
//   - *int32: A pointer to the int32 value if valid, or nil if the NullInt32 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt32{Int32: 42, Valid: true}
//	ptr := Int32Nil(nullInt) // Returns: &42
//
//	invalidInt := sql.NullInt32{Int32: 0, Valid: false}
//	ptr := Int32Nil(invalidInt) // Returns: nil
func Int32Nil(nullInt sql.NullInt32) *int32 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int32
}

// Int32 creates a sql.NullInt32 from type-constrained input types.
// This function provides a type-safe way to create nullable int32 values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - int32: Creates a valid NullInt32 with the given int32 value
//   - *int32: Creates a valid NullInt32 from pointer (nil pointer creates invalid NullInt32)
//
// Returns:
//   - sql.NullInt32: A NullInt32 struct with appropriate Valid flag and Int32 value.
//
// Examples:
//
//	// From int32 value
//	nullInt := Int32(int32(42)) // Returns: {Int32: 42, Valid: true}
//
//	// From int32 pointer
//	intPtr := int32(100)
//	nullInt := Int32(&intPtr) // Returns: {Int32: 100, Valid: true}
//
//	// From nil pointer
//	nullInt := Int32((*int32)(nil)) // Returns: {Int32: 0, Valid: false}
func Int32[T int32Type](val T) sql.NullInt32 {
	switch v := any(val).(type) {
	case int32:
		return sql.NullInt32{
			Int32: v,
			Valid: true,
		}
	case *int32:
		if v == nil {
			return sql.NullInt32{
				Int32: 0,
				Valid: false,
			}
		}
		return sql.NullInt32{
			Int32: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt32
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}
}

// Int16Any converts a sql.NullInt16 to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable int16 value
// to database driver operations.
//
// Parameters:
//   - nullInt (sql.NullInt16): The nullable int16 value to convert.
//
// Returns:
//   - driver.Value: The int16 value if valid, or nil if the NullInt16 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt16{Int16: 42, Valid: true}
//	value := Int16Any(nullInt) // Returns: 42
//
//	invalidInt := sql.NullInt16{Int16: 0, Valid: false}
//	value := Int16Any(invalidInt) // Returns: nil
func Int16Any(nullInt sql.NullInt16) driver.Value {
	if !nullInt.Valid {
		return nil
	}
	return nullInt.Int16
}

// Int16Nil converts a sql.NullInt16 to a pointer to int16 (*int16).
// This function is useful when you need to work with nullable int16 values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullInt (sql.NullInt16): The nullable int16 value to convert.
//
// Returns:
//   - *int16: A pointer to the int16 value if valid, or nil if the NullInt16 is invalid/null.
//
// Example:
//
//	nullInt := sql.NullInt16{Int16: 42, Valid: true}
//	ptr := Int16Nil(nullInt) // Returns: &42
//
//	invalidInt := sql.NullInt16{Int16: 0, Valid: false}
//	ptr := Int16Nil(invalidInt) // Returns: nil
func Int16Nil(nullInt sql.NullInt16) *int16 {
	if !nullInt.Valid {
		return nil
	}
	return &nullInt.Int16
}

// Int16 creates a sql.NullInt16 from type-constrained input types.
// This function provides a type-safe way to create nullable int16 values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - int16: Creates a valid NullInt16 with the given int16 value
//   - *int16: Creates a valid NullInt16 from pointer (nil pointer creates invalid NullInt16)
//
// Returns:
//   - sql.NullInt16: A NullInt16 struct with appropriate Valid flag and Int16 value.
//
// Examples:
//
//	// From int16 value
//	nullInt := Int16(int16(42)) // Returns: {Int16: 42, Valid: true}
//
//	// From int16 pointer
//	intPtr := int16(100)
//	nullInt := Int16(&intPtr) // Returns: {Int16: 100, Valid: true}
//
//	// From nil pointer
//	nullInt := Int16((*int16)(nil)) // Returns: {Int16: 0, Valid: false}
func Int16[T int16Type](val T) sql.NullInt16 {
	switch v := any(val).(type) {
	case int16:
		return sql.NullInt16{
			Int16: v,
			Valid: true,
		}
	case *int16:
		if v == nil {
			return sql.NullInt16{
				Int16: 0,
				Valid: false,
			}
		}
		return sql.NullInt16{
			Int16: *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullInt16
		return sql.NullInt16{
			Int16: 0,
			Valid: false,
		}
	}
}
