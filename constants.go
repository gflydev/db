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
