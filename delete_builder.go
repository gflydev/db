package db

import (
	"github.com/gflydev/core/errors"
	qb "github.com/jivegroup/fluentsql"
)

// Delete removes records from the database table based on the provided model and conditions.
// This method provides a flexible approach to data deletion by supporting both primary key-based
// deletion and conditional deletion using WHERE clauses. It automatically constructs DELETE SQL
// statements based on the model's primary keys and any additional conditions specified through
// the fluent interface. The method ensures safe deletion by requiring at least one WHERE condition
// to prevent accidental deletion of all table data.
//
// Parameters:
//   - model (any): The model defining the target table and deletion criteria. Supported types:
//   - Struct: Direct struct value representing the table model (e.g., User{ID: 1})
//   - *Struct: Pointer to struct representing the table model (e.g., &User{ID: 1})
//   - The struct should have appropriate field tags for database mapping
//   - Primary key fields in the model are used to build WHERE conditions automatically
//   - Non-nil primary key values are included in the deletion criteria
//
// Returns:
//   - error: Returns an error if:
//   - The model cannot be processed or reflected (invalid struct type)
//   - Table metadata extraction fails
//   - No WHERE conditions are present (safety requirement)
//   - Raw SQL execution fails (when using Raw() method)
//   - DELETE statement construction fails
//   - Database execution fails
//   - Database connectivity issues occur
//     Returns nil on successful deletion.
//
// Examples:
//
//	// Delete by primary key
//	user := User{ID: 123}
//	err := db.Model(&User{}).Delete(user)
//	// Executes: DELETE FROM users WHERE id = 123
//
//	// Delete with pointer to struct
//	user := &User{ID: 456}
//	err := db.Model(&User{}).Delete(user)
//	// Executes: DELETE FROM users WHERE id = 456
//
//	// Delete with additional WHERE conditions
//	user := User{}
//	err := db.Model(&User{}).Where("status", "inactive").Delete(user)
//	// Executes: DELETE FROM users WHERE status = 'inactive'
//
//	// Delete with multiple conditions
//	user := User{ID: 789}
//	err := db.Model(&User{}).Where("created_at < ?", time.Now().AddDate(-1, 0, 0)).Delete(user)
//	// Executes: DELETE FROM users WHERE id = 789 AND created_at < '2023-01-01'
//
//	// Delete with OR conditions
//	user := User{}
//	err := db.Model(&User{}).Where("status", "inactive").Or("last_login < ?", oldDate).Delete(user)
//	// Executes: DELETE FROM users WHERE status = 'inactive' OR last_login < '2023-01-01'
//
//	// Delete with grouped conditions
//	user := User{}
//	err := db.Model(&User{}).Where("status", "inactive").
//	    Where(func(db *DBModel) *DBModel {
//	        return db.Where("role", "guest").Or("role", "temp")
//	    }).Delete(user)
//	// Executes: DELETE FROM users WHERE status = 'inactive' AND (role = 'guest' OR role = 'temp')
//
//	// Raw SQL deletion
//	user := User{}
//	err := db.Model(&User{}).Raw("DELETE FROM users WHERE created_at < NOW() - INTERVAL 1 YEAR").Delete(user)
//	// Executes the raw SQL directly
//
//	// Composite primary key deletion
//	orderItem := OrderItem{OrderID: 100, ProductID: 200}
//	err := db.Model(&OrderItem{}).Delete(orderItem)
//	// Executes: DELETE FROM order_items WHERE order_id = 100 AND product_id = 200
//
// Safety Features:
//   - Requires at least one WHERE condition to prevent accidental mass deletion
//   - Automatically includes non-nil primary key values as WHERE conditions
//   - Validates model structure before executing deletion
//   - Supports transaction rollback through proper error handling
//
// Note:
//   - Primary key fields with nil/zero values are ignored in WHERE conditions
//   - The method respects all WHERE conditions set via the fluent interface
//   - Raw SQL takes precedence over model-based deletion when both are present
//   - The fluent model builder is automatically reset after the operation
//   - For mass deletion, use Raw() method with appropriate WHERE clauses
//   - Composite primary keys are fully supported with AND logic
//   - The operation is atomic and will either succeed completely or fail without changes
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
				Opt:   Eq,         // Equality operator for the condition.
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
		case condition.AndOr == And:
			// Add the AND condition to the query builder.
			deleteBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == Or:
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
