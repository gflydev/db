package db

// Where represents a collection of conditions that form a SQL WHERE clause.
// It encapsulates multiple Condition structs that can be combined to create
// complex filtering criteria.
type Where struct {
	// Conditions represent a slice of Condition structs that define the WHERE clause of a SQL query.
	// Each Condition specifies filtering criteria like field comparisons and logical operations.
	Conditions []Condition
}

// WhereBuilder provides a fluent interface for constructing SQL WHERE clauses.
// It wraps a Where struct and provides methods to incrementally build conditions
// through method chaining.
type WhereBuilder struct {
	whereStatement Where // whereStatement holds the WHERE conditions of the query.
}

// FnWhereBuilder represents a function type that takes a WhereBuilder and returns a WhereBuilder pointer.
// This function type is used to build grouped conditions by allowing custom condition building logic
// to be passed as a parameter.
//
// Example:
//
//	whereBuilder.WhereGroup(func(wb WhereBuilder) *WhereBuilder {
//		return wb.Where("status", "active").WhereOr("type", "premium")
//	})
type FnWhereBuilder func(whereBuilder WhereBuilder) *WhereBuilder

// WhereInstance creates a new WhereBuilder instance.
// This function provides a constructor for WhereBuilder similar to qb.WhereInstance().
//
// Returns:
//   - *WhereBuilder: A new WhereBuilder instance with empty conditions.
func WhereInstance() *WhereBuilder {
	return &WhereBuilder{
		whereStatement: Where{
			Conditions: []Condition{},
		},
	}
}

// Conditions returns the slice of conditions from the WhereBuilder.
// This method provides access to the internal conditions similar to qb.WhereBuilder.Conditions().
//
// Returns:
//   - []Condition: The slice of conditions in the WhereBuilder.
func (wb *WhereBuilder) Conditions() []Condition {
	return wb.whereStatement.Conditions
}

// Where adds a WHERE condition to the WhereBuilder.
// This method allows building conditions within the WhereBuilder.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *WhereBuilder: A reference to the WhereBuilder instance for chaining.
func (wb *WhereBuilder) Where(field any, opt WhereOpt, value any) *WhereBuilder {
	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
	})
	return wb
}

// WhereOr adds an OR WHERE condition to the WhereBuilder.
// This method allows building OR conditions within the WhereBuilder.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *WhereBuilder: A reference to the WhereBuilder instance for chaining.
func (wb *WhereBuilder) WhereOr(field any, opt WhereOpt, value any) *WhereBuilder {
	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: Or,
	})
	return wb
}

// WhereGroup adds a grouped WHERE condition to the WhereBuilder.
// This method allows building grouped conditions (enclosed in parentheses) within the WhereBuilder.
//
// Parameters:
//   - groupCondition (FnWhereBuilder): A function that builds the conditions within the group.
//
// Returns:
//   - *WhereBuilder: A reference to the WhereBuilder instance for chaining.
func (wb *WhereBuilder) WhereGroup(groupCondition FnWhereBuilder) *WhereBuilder {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*WhereInstance())

	cond := Condition{
		Group: whereBuilder.whereStatement.Conditions,
	}

	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, cond)

	return wb
}

// WhereCondition adds multiple WHERE conditions to the WhereBuilder.
// This method allows adding pre-built conditions directly to the WhereBuilder.
//
// Parameters:
//   - conditions (...Condition): Variable number of Condition structs to add.
//
// Returns:
//   - *WhereBuilder: A reference to the WhereBuilder instance for chaining.
func (wb *WhereBuilder) WhereCondition(conditions ...Condition) *WhereBuilder {
	wb.whereStatement.Conditions = append(wb.whereStatement.Conditions, conditions...)
	return wb
}
