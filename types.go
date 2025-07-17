package db

import (
	qb "github.com/jivegroup/fluentsql"
)

// ====================================================================
//                         Join Type
// ====================================================================

// JoinType represents the type of SQL join operation.
// This is an alias for fluentsql.JoinType to maintain backward compatibility.
type JoinType = qb.JoinType

// Join type constants from fluentsql package
const (
	InnerJoin     = qb.InnerJoin     // Inner join
	LeftJoin      = qb.LeftJoin      // Left outer join
	RightJoin     = qb.RightJoin     // Right outer join
	FullOuterJoin = qb.FullOuterJoin // Full outer join
	CrossJoin     = qb.CrossJoin     // Cross join
)

// ====================================================================
//                         Order By Type
// ====================================================================

// OrderByDir represents the sorting direction.
// This is an alias for fluentsql.OrderByDir to maintain backward compatibility.
//
// Values:
// - Asc: Ascending order.
// - Desc: Descending order.
type OrderByDir = qb.OrderByDir

// Constants representing sorting directions from fluentsql package
const (
	Asc  = qb.Asc  // Ascending order.
	Desc = qb.Desc // Descending order.
)

// ====================================================================
//                         Where And Or Type
// ====================================================================

// WhereAndOr represents logical operators for combining conditions.
// This is an alias for fluentsql.WhereAndOr to maintain backward compatibility.
type WhereAndOr = qb.WhereAndOr

// Logical operator constants from fluentsql package
const (
	And = qb.And // Logical AND operator for combining conditions
	Or  = qb.Or  // Logical OR operator for combining conditions
)

// WhereOpt defines the operators used in SQL conditions.
// This is an alias for fluentsql.WhereOpt to maintain backward compatibility.
type WhereOpt = qb.WhereOpt

// SQL condition operator constants from fluentsql package
const (
	Eq         = qb.Eq         // Equal to (=)
	NotEq      = qb.NotEq      // Not equal to (<>)
	Diff       = qb.Diff       // Not equal to (!=)
	Greater    = qb.Greater    // Greater than (>)
	Lesser     = qb.Lesser     // Less than (<)
	GrEq       = qb.GrEq       // Greater than or equal to (>=)
	LeEq       = qb.LeEq       // Less than or equal to (<=)
	Like       = qb.Like       // Pattern matching (LIKE)
	NotLike    = qb.NotLike    // Not pattern matching (NOT LIKE)
	In         = qb.In         // Value in a list (IN)
	NotIn      = qb.NotIn      // Value not in a list (NOT IN)
	Between    = qb.Between    // Value in a range (BETWEEN)
	NotBetween = qb.NotBetween // Value not in a range (NOT BETWEEN)
	Null       = qb.Null       // Null value (IS NULL)
	NotNull    = qb.NotNull    // Not null value (IS NOT NULL)
	Exists     = qb.Exists     // Subquery results exist (EXISTS)
	NotExists  = qb.NotExists  // Subquery results do not exist (NOT EXISTS)
	EqAny      = qb.EqAny      // Equal to any value in a subquery (= ANY)
	NotEqAny   = qb.NotEqAny   // Not equal to any value in a subquery (<> ANY)
	DiffAny    = qb.DiffAny    // Not equal to any value in a subquery (!= ANY)
	GreaterAny = qb.GreaterAny // Greater than any value in a subquery (> ANY)
	LesserAny  = qb.LesserAny  // Less than any value in a subquery (< ANY)
	GrEqAny    = qb.GrEqAny    // Greater than or equal to any value in a subquery (>= ANY)
	LeEqAny    = qb.LeEqAny    // Less than or equal to any value in a subquery (<= ANY)
	EqAll      = qb.EqAll      // Equal to all values in a subquery (= ALL)
	NotEqAll   = qb.NotEqAll   // Not equal to all values in a subquery (<> ALL)
	DiffAll    = qb.DiffAll    // Not equal to all values in a subquery (!= ALL)
	GreaterAll = qb.GreaterAll // Greater than all values in a subquery (> ALL)
	LesserAll  = qb.LesserAll  // Less than all values in a subquery (< ALL)
	GrEqAll    = qb.GrEqAll    // Greater than or equal to all values in a subquery (>= ALL)
	LeEqAll    = qb.LeEqAll    // Less than or equal to all values in a subquery (<= ALL)
)

// ====================================================================
//                         Other Types
// ====================================================================

// ValueField represents a column/field in a SQL query as a string value.
// This is an alias for fluentsql.ValueField to maintain backward compatibility.
//
// Usage:
//   - Used to reference database columns in SQL queries
//   - Provides type safety for field references
//
// Methods:
//   - String(): Converts the ValueField to its string representation.
//
// Example:
//
//	qb.ValueField("user_details.user_id")
type ValueField qb.ValueField

// Make sure ValueField implements the IValueField interface.
var _ qb.IValueField = (*ValueField)(nil)

// Value returns the string representation of the ValueField.
// This is used to convert the ValueField type to a plain string
// for use in SQL queries and string operations.
//
// Returns:
//   - string: The string value of the ValueField
//
// Example:
//
//	field := ValueField("user.id")
//	str := field.String() // Returns "user.id"
func (v ValueField) Value() string {
	return string(v)
}

// Limit represents the SQL LIMIT and OFFSET clauses for pagination.
//
// Usage:
//   - Controls the number and starting point of returned rows
//   - Used for implementing pagination in database queries
//
// Example:
//
//	LIMIT 10 OFFSET 20
type Limit qb.Limit

// Fetch clause represents a SQL FETCH clause with offset and limit.
// This is an alias for fluentsql.Fetch to maintain backward compatibility.
//
// Usage:
//   - Controls pagination through OFFSET and FETCH NEXT clauses
//   - Supported in MSSQL Server and Oracle databases
//
// Example:
//
//	OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY
type Fetch qb.Fetch

// FieldNot represents a SQL field prefixed with a NOT operator for negating conditions.
// This is an alias for fluentsql.FieldNot to maintain backward compatibility.
//
// Usage:
//   - Used to create negated field conditions in SQL queries
//   - Applies NOT operator to field expressions
//
// Methods:
//   - String(): Returns the SQL string representation of the negated field.
type FieldNot qb.FieldNot

// FieldEmpty represents an empty SQL field, often used in conditions like EXISTS or NOT EXISTS.
// This is an alias for fluentsql.FieldEmpty to maintain backward compatibility.
//
// Usage:
//   - Used in subquery conditions where no specific field is referenced
//   - Commonly used with EXISTS and NOT EXISTS operations
//
// Example:
//
//	WHERE NOT EXISTS (SELECT employee_id FROM dependents)
type FieldEmpty qb.FieldEmpty

// FieldYear represents a SQL year extraction operation for a given field.
// This is an alias for fluentsql.FieldYear to maintain backward compatibility.
//
// Usage:
//   - Extracts the year portion from date/datetime fields
//   - Generates database-specific SQL syntax for year extraction
//
// Database-specific implementations:
//   - MySQL: YEAR(hire_date) Between 1990 AND 1993
//   - PostgreSQL: DATE_PART('year', hire_date) Between 1990 AND 1993
//   - SQLite: strftime('%Y', hire_date)
type FieldYear qb.FieldYear

// Condition represents a single condition in a WHERE clause.
// This is an alias for fluentsql.Condition to maintain backward compatibility.
//
// Usage:
//   - Defines comparison operations between fields and values
//   - Supports standard SQL operators (=, >, <, LIKE, etc.)
//   - Can be combined to form complex conditions
//
// Example:
//
//	WHERE user_id = 1 AND status = 'active'
//
// Condition type struct
type Condition struct {
	// Field represents the name of the column to compare. It can be of type `string` or `FieldNot`.
	Field any
	// Opt specifies the condition operator such as =, <>, >, <, >=, <=, LIKE, IN, NOT IN, BETWEEN, etc.
	Opt WhereOpt
	// Value holds the value to be compared against the field. Support ValueField for checking with table's column
	Value any
	// AndOr specifies the logical combination with the previous condition (AND, OR). Default is AND.
	AndOr WhereAndOr
	// Group contains sub-conditions enclosed in parentheses `()`.
	Group []Condition
}

// ToQBCondition converts a Condition to qb.Condition.
// This helper function recursively converts the local Condition struct
// to the fluentsql qb.Condition struct, including any nested Group conditions.
//
// Returns:
//   - qb.Condition: The converted condition ready for use with fluentsql
//
// Example:
//
//	condition := Condition{
//		Field: "user_id",
//		Opt:   Eq,
//		Value: 123,
//		AndOr: And,
//	}
//	qbCondition := condition.ToQBCondition()
func (c Condition) ToQBCondition() qb.Condition {
	// Convert nested Group conditions recursively
	var qbGroup []qb.Condition
	if len(c.Group) > 0 {
		qbGroup = make([]qb.Condition, len(c.Group))
		for i, groupCondition := range c.Group {
			qbGroup[i] = groupCondition.ToQBCondition()
		}
	}

	return qb.Condition{
		Field: c.Field,
		Opt:   c.Opt,
		Value: c.Value,
		AndOr: c.AndOr,
		Group: qbGroup,
	}
}
