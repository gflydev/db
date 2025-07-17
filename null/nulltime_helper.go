package null

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// timeConstraint is a constraint interface that allows either time.Time or *time.Time types.
// It is used in generic functions to handle both direct time.Time values and pointers
// to time.Time values in a type-safe manner.
type timeConstraint interface {
	time.Time | *time.Time
}

// TimeAny converts a sql.NullTime to a driver.Value for database operations.
// This function is typically used when you need to pass a nullable time.Time value
// to database driver operations.
//
// Parameters:
//   - nullTime (sql.NullTime): The nullable time.Time value to convert.
//
// Returns:
//   - driver.Value: The time.Time value if valid, or nil if the NullTime is invalid/null.
//
// Example:
//
//	now := time.Now()
//	nullTime := sql.NullTime{Time: now, Valid: true}
//	value := TimeAny(nullTime) // Returns: now
//
//	invalidTime := sql.NullTime{Time: time.Time{}, Valid: false}
//	value := TimeAny(invalidTime) // Returns: nil
func TimeAny(nullTime sql.NullTime) driver.Value {
	if !nullTime.Valid {
		return nil
	}
	return nullTime.Time
}

// TimeNil converts a sql.NullTime to a pointer to time.Time (*time.Time).
// This function is useful when you need to work with nullable time.Time values
// in your application logic, where nil represents a null database value.
//
// Parameters:
//   - nullTime (sql.NullTime): The nullable time.Time value to convert.
//
// Returns:
//   - *time.Time: A pointer to the time.Time value if valid, or nil if the NullTime is invalid/null.
//
// Example:
//
//	now := time.Now()
//	nullTime := sql.NullTime{Time: now, Valid: true}
//	ptr := TimeNil(nullTime) // Returns: &now
//
//	invalidTime := sql.NullTime{Time: time.Time{}, Valid: false}
//	ptr := TimeNil(invalidTime) // Returns: nil
func TimeNil(nullTime sql.NullTime) *time.Time {
	if !nullTime.Valid {
		return nil
	}
	return &nullTime.Time
}

// Time creates a sql.NullTime from type-constrained input types.
// This function provides a type-safe way to create nullable time.Time values
// for database operations, handling both direct values and pointers with compile-time type checking.
//
// Parameters:
//   - val (T): The input value to convert. Supported types:
//   - time.Time: Creates a valid NullTime with the given time.Time value
//   - *time.Time: Creates a valid NullTime from pointer (nil pointer creates invalid NullTime)
//
// Returns:
//   - sql.NullTime: A NullTime struct with appropriate Valid flag and Time value.
//
// Examples:
//
//	// From time.Time value
//	now := time.Now()
//	nullTime := Time(now) // Returns: {Time: now, Valid: true}
//
//	// From time.Time pointer
//	timePtr := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//	nullTime := Time(&timePtr) // Returns: {Time: timePtr, Valid: true}
//
//	// From nil pointer
//	nullTime := Time((*time.Time)(nil)) // Returns: {Time: time.Time{}, Valid: false}
func Time[T timeConstraint](val T) sql.NullTime {
	switch v := any(val).(type) {
	case time.Time:
		return sql.NullTime{
			Time:  v,
			Valid: true,
		}
	case *time.Time:
		if v == nil {
			return sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			}
		}
		return sql.NullTime{
			Time:  *v,
			Valid: true,
		}
	default:
		// For any other type, return invalid NullTime
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
}

// TimeNow creates a sql.NullTime with the current time.
// This is a convenience function that creates a valid NullTime using the current
// system time, equivalent to calling Time(time.Now()).
//
// Returns:
//   - sql.NullTime: A valid NullTime struct with the current time and Valid set to true.
//
// Example:
//
//	nullTime := TimeNow() // Returns: {Time: <current time>, Valid: true}
//
//	// Equivalent to:
//	nullTime := Time(time.Now())
func TimeNow() sql.NullTime {
	return Time(time.Now())
}
