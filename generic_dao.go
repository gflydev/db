package db

import (
	"database/sql"
	"errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	qb "github.com/jiveio/fluentsql" // Query builder
)

// ====================================================================
// ========================= Specific methods =========================
// ====================================================================

// GetModelByID retrieves the first record of type T from the database
// where the "id" field matches the specified value.
//
// This is a convenience wrapper around GetModelBy for querying records by their unique identifier.
//
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - value (any): The value of the "id" field to match against.
//
// Returns:
//   - *T: A pointer to the retrieved model of type T, or nil if no matching record is found.
//   - error: An error object if an error occurs during the retrieval process.
func GetModelByID[T any](value any) (*T, error) {
	return GetModelBy[T]("id", value)
}

// GetModelBy allows filtering records of type T from the database
// by specifying a field and its required value.
//
// This is a helper function that makes use of the GetModelWhereEq function to apply an equality condition.
//
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - field (string): The name of the database field to filter on.
//   - value (any): The value the specified field is required to equal.
//
// Returns:
//   - *T: A pointer to the first matching record of type T retrieved from the database, or nil if no record is found.
//   - error: An error object if an error occurs during the retrieval process.
func GetModelBy[T any](field string, value any) (*T, error) {
	return GetModelWhereEq[T](field, value)
}

// ====================================================================
// ========================== Generic methods =========================
// ====================================================================

// GetModelWhereEq retrieves the first record of type T from the database
// where the specified field equals the given value.
//
// Parameters:
//   - m: A pointer to the model where the result will be stored.
//   - field: The name of the database field to filter by.
//   - value: The value to match the field against.
//
// Returns:
//   - error: An error object if an error occurs during the retrieval process.
//     Returns nil if the query succeeds.
func GetModelWhereEq[T any](field string, value any) (*T, error) {
	return GetModel[T](qb.Condition{
		Field: field,
		Opt:   qb.Eq,
		Value: value,
	})
}

// GetModel retrieves the first record of type T from the database
// that matches the provided conditions.
//
// Parameters:
//   - m: A pointer to the model where the result will be stored.
//   - conditions: Variadic list of qb.Condition specifying the field, operator,
//     and value to filter the query.
//
// Returns:
//   - error: An error object if an error occurs during the retrieval process.
//     Returns nil if the query succeeds. Logs unexpected errors.
func GetModel[T any](conditions ...qb.Condition) (*T, error) {
	var err error
	var m T

	// Try/catch block
	try.Perform(func() {
		builder := Instance()
		for _, condition := range conditions {
			builder.Where(condition.Field, condition.Opt, condition.Value)
		}

		// Get first record then assign to `m` WHERE field = value

		if e := builder.First(&m); e != nil {
			try.Throw(e)
		}
	}).Catch(func(e try.E) {
		err = e.(error)

		// Log unexpected error!
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error(e)
		}
	})

	return &m, err
}

// FindModels retrieves a paginated list of records of type T from the database
// that match the provided conditions.
//
// Parameters:
//   - page (int): The current page number (1-based). Defaults to 0 if not provided.
//   - limit (int): The number of records to retrieve per page.
//   - conditions (...qb.Condition): Variadic list of qb.Condition specifying the field,
//     operator, and value to filter the query.
//
// Returns:
//   - ([]T): A slice of records of type T.
//   - (int): The total number of records that match the conditions.
//   - (error): An error object if an error occurs during the retrieval process.
func FindModels[T any](page, limit int, sortField string, sortDir qb.OrderByDir, conditions ...qb.Condition) ([]T, int, error) {
	var items []T
	var total int
	var err error

	var offset = 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	try.Perform(func() {
		builder := Instance()
		for _, condition := range conditions {
			builder.Where(condition.Field, condition.Opt, condition.Value)
		}

		builder.OrderBy(sortField, sortDir)

		total, err = builder.Limit(limit, offset).Find(&items)
		if err != nil {
			try.Throw(err)
		}
	}).Catch(func(e try.E) {
		err = e.(error)

		// Log unexpected error!
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error(e)
		}
	})

	// For case empty list => return an empty []T
	if items == nil || err != nil {
		items = []T{}
	}

	return items, total, err
}

// CreateModel creates a new record of type T in the database.
// It begins a transaction, attempts to create the record, and commits the transaction.
// If an error occurs, the transaction is rolled back and the error is returned.
//
// Parameters:
//   - m: A pointer to the model to be created.
//
// Returns:
//   - error: An error object if an error occurs during the creation process.
func CreateModel[T any](m *T) error {
	var err error
	db := Instance()

	try.Perform(func() {
		// Begin transaction
		db.Begin()
		// Trying to create an instance
		if e := db.Create(m); e != nil {
			try.Throw(e)
		}
		// Commit transaction
		if e := db.Commit(); e != nil {
			try.Throw(e)
		}
	}).Catch(func(e try.E) {
		err = e.(error)
		// Rollback transaction
		_ = db.Rollback()
	})

	return err
}

// UpdateModel updates a record of type T in the database.
// It begins a transaction, updates the record, and commits the transaction.
// If an error occurs, the transaction is rolled back and the error is returned.
//
// Parameters:
//   - m: A pointer to the model to be updated.
//
// Returns:
//   - error: An error object if an error occurs during the update process.
func UpdateModel[T any](m *T) error {
	var err error
	db := Instance()

	try.Perform(func() {
		// Begin transaction
		db.Begin()

		// Attempt to update the record
		if e := db.Update(m); e != nil {
			try.Throw(e)
		}

		// Commit transaction
		if e := db.Commit(); e != nil {
			try.Throw(e)
		}
	}).Catch(func(e try.E) {
		// Handle error and rollback transaction
		err = e.(error)
		_ = db.Rollback()
	})

	return err
}

// DeleteModel deletes a record of type T from the database.
// It begins a transaction, deletes the record, and commits the transaction.
// If an error occurs, the transaction is rolled back and the error is returned.
//
// Parameters:
//   - m: A pointer to the model to be deleted.
//
// Returns:
//   - error: An error object if an error occurs during the deletion process.
func DeleteModel[T any](m *T) error {
	var err error
	db := Instance()

	try.Perform(func() {
		// Begin transaction
		db.Begin()

		// Attempt to delete the record
		if e := db.Delete(m); e != nil {
			try.Throw(e)
		}

		// Commit transaction
		if e := db.Commit(); e != nil {
			try.Throw(e)
		}
	}).Catch(func(e try.E) {
		// Handle error and rollback transaction
		err = e.(error)
		_ = db.Rollback()
	})

	return err
}
