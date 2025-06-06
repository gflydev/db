package db

import (
	"github.com/gflydev/core/errors"
	qb "github.com/jivegroup/fluentsql"
)

// Delete performs the deletion of data for a given table using a model of type Struct or *Struct.
//
// Parameters:
//   - model (any): The input model defining the data to delete. This can be a struct or a pointer to a struct.
//
// Returns:
//   - error: Returns an error if the deletion process fails.
func (db *DBModel) Delete(model any) error {
	var err error // Stores errors encountered during the function execution.

	// Delete using raw SQL if it's set.
	if db.raw.sqlStr != "" {
		err = db.execRaw(db.raw.sqlStr, db.raw.args)
		if err != nil {
			return err
		}

		// Reset fluent model builder.
		db.reset()
	}

	var table *Table         // Represents the table corresponding to the model.
	var hasCondition = false // Indicates if any WHERE condition is present.

	// Create a table object from the given model.
	if table, err = ModelData(model); err != nil {
		return err
	}

	// Create an instance of a delete query builder.
	deleteBuilder := qb.DeleteInstance().Delete(table.Name)

	// Build WHERE clause using primary columns of the table.
	for _, primaryColumn := range table.Primaries {
		primaryKey := primaryColumn.Name       // The name of the primary column.
		primaryVal := table.Values[primaryKey] // The value of the primary column.

		if primaryVal != nil {
			wherePrimaryCondition := qb.Condition{
				Field: primaryKey, // Field name for the condition.
				Opt:   qb.Eq,      // Equality operator for the condition.
				Value: primaryVal, // Value to match against.
				AndOr: qb.And,     // Logical operator for chaining conditions.
			}

			// Add the primary key condition to the query builder.
			deleteBuilder.WhereCondition(wherePrimaryCondition)
			hasCondition = true // Mark that at least one condition is present.
		}
	}

	// Build WHERE clause using additional conditions from the condition list.
	for _, condition := range db.whereStatement.Conditions {
		switch {
		case len(condition.Group) > 0:
			// Append grouped conditions to the query builder.
			deleteBuilder.WhereGroup(func(whereBuilder qb.WhereBuilder) *qb.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)
				return &whereBuilder
			})
			hasCondition = true
		case condition.AndOr == qb.And:
			// Add the AND condition to the query builder.
			deleteBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == qb.Or:
			// Add the OR condition to the query builder.
			deleteBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		}
	}

	// Ensure there is at least one WHERE condition.
	if !hasCondition {
		return errors.New("Missing WHERE condition for deleting operator")
	}

	// Execute the delete operation using the constructed delete builder.
	err = db.delete(deleteBuilder)

	// Reset fluent model builder.
	db.reset()

	return err
}
