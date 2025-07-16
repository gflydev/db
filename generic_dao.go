package db

import (
	"database/sql"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	qb "github.com/jivegroup/fluentsql" // Query builder
)

// ====================================================================
//                           Specific methods
// ====================================================================

// GetModelByID retrieves a single database record by its primary key identifier.
// This is the most commonly used function for fetching individual records when you know
// the primary key value. It provides a convenient, type-safe way to retrieve records
// without manually constructing WHERE conditions. The function assumes the primary key
// field is named "id" and automatically handles type conversion and error handling.
//
// Generic Type:
//   - T: The model struct type that represents the database table. Must have appropriate
//     database tags for column mapping and should include an "id" field as the primary key.
//
// Parameters:
//   - value (any): The primary key value to search for. Supported types:
//   - int, int64, uint, uint64: For integer primary keys
//   - string: For string-based primary keys (UUIDs, etc.)
//   - Other types that can be converted to database-compatible values
//     The value will be automatically converted to match the target field type.
//
// Returns:
//   - *T: A pointer to the retrieved model instance with all fields populated from the database.
//     Returns nil if no record is found with the specified ID.
//   - error: Returns an error if:
//   - Database connection fails
//   - The model type T cannot be processed
//   - SQL execution fails
//   - Type conversion fails
//     Returns errors.ItemNotFound if no record exists with the given ID.
//     Returns nil on successful retrieval.
//
// Examples:
//
//	// Retrieve user by integer ID
//	type User struct {
//	    ID    int64  `db:"id,primary"`
//	    Name  string `db:"name"`
//	    Email string `db:"email"`
//	}
//
//	user, err := GetModelByID[User](123)
//	if err != nil {
//	    if errors.Is(err, errors.ItemNotFound) {
//	        log.Println("User not found")
//	    } else {
//	        log.Printf("Database error: %v", err)
//	    }
//	    return
//	}
//	fmt.Printf("Found user: %s (%s)", user.Name, user.Email)
//
//	// Retrieve product by string ID (UUID)
//	type Product struct {
//	    ID    string  `db:"id,primary"`
//	    Name  string  `db:"name"`
//	    Price float64 `db:"price"`
//	}
//
//	productID := "550e8400-e29b-41d4-a716-446655440000"
//	product, err := GetModelByID[Product](productID)
//	if err == nil {
//	    fmt.Printf("Product: %s - $%.2f", product.Name, product.Price)
//	}
//
//	// Error handling patterns
//	order, err := GetModelByID[Order](456)
//	switch {
//	case errors.Is(err, errors.ItemNotFound):
//	    // Handle not found case
//	    return nil, fmt.Errorf("order %d does not exist", 456)
//	case err != nil:
//	    // Handle other database errors
//	    return nil, fmt.Errorf("failed to retrieve order: %w", err)
//	default:
//	    // Success case
//	    return order, nil
//	}
//
//	// Batch retrieval with error handling
//	userIDs := []int64{1, 2, 3, 4, 5}
//	var users []*User
//	for _, id := range userIDs {
//	    user, err := GetModelByID[User](id)
//	    if err != nil && !errors.Is(err, errors.ItemNotFound) {
//	        return nil, err // Stop on database errors
//	    }
//	    if user != nil {
//	        users = append(users, user)
//	    }
//	}
//
// Use Cases:
//   - Retrieving user profiles by user ID
//   - Fetching specific orders, products, or entities by their primary key
//   - Loading related records when you have foreign key references
//   - Implementing REST API endpoints that accept ID parameters
//   - Cache lookups where you need to fetch missing records by ID
//
// Performance Notes:
//   - Uses primary key index for optimal query performance
//   - Single database round trip for each call
//   - Efficient for individual record retrieval
//   - Consider batch operations for multiple records
//
// Note:
//   - Assumes the primary key field is named "id"
//   - Returns errors.ItemNotFound for missing records (not sql.ErrNoRows)
//   - The model type T must have proper database tags
//   - Thread-safe and can be called concurrently
//   - Automatically handles database connection management
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
	item, err := GetModelWhereEq[T](field, value)
	// Log unexpected error!
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.ItemNotFound
	}

	return item, err
}

// ====================================================================
//                           Generic methods
// ====================================================================

// GetModelWhereEq retrieves the first record of type T from the database
// where the specified field equals the given value.
//
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - field (string): The name of the database field to filter by.
//   - value (any): The value to match the field against.
//
// Returns:
//   - *T: A pointer to the retrieved model of type T, or nil if no matching record is found.
//   - error: An error object if an error occurs during the retrieval process.
func GetModelWhereEq[T any](field string, value any) (*T, error) {
	return GetModel[T](Condition{
		Field: field,
		Opt:   Eq,
		Value: value,
	})
}

// GetModel retrieves the first record of type T from the database
// that matches the provided conditions.
//
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - conditions (...qb.Condition): Variadic list of qb.Condition specifying the field, operator,
//     and value to filter the query.
//
// Returns:
//   - *T: A pointer to the retrieved model of type T, or nil if no matching record is found.
//   - error: An error object if an error occurs during the retrieval process.
//     Returns nil if the query succeeds. Logs unexpected errors.
func GetModel[T any](conditions ...Condition) (*T, error) {
	var builder = Instance()
	var err error
	var m T

	// Try/catch block
	try.Perform(func() {
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
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.ItemNotFound
		} else {
			log.Error(e)
		}
	})

	return &m, err
}

// FindModels retrieves a paginated list of records of type T from the database
// that match the provided conditions.
//
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - page (int): The current page number (1-based). Defaults to 0 if not provided.
//   - limit (int): The number of records to retrieve per page.
//   - sortField (string): The field name to sort the results by.
//   - sortDir (qb.OrderByDir): The sorting direction (qb.Asc for ascending, qb.Desc for descending).
//   - conditions (...qb.Condition): Variadic list of qb.Condition specifying the field,
//     operator, and value to filter the query.
//
// Returns:
//   - []T: A slice of records of type T.
//   - int: The total number of records that match the conditions.
//   - error: An error object if an error occurs during the retrieval process.
func FindModels[T any](page, limit int, sortField string, sortDir qb.OrderByDir, conditions ...qb.Condition) ([]T, int, error) {
	var builder = Instance()
	var items []T
	var total int
	var err error

	var offset = 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	try.Perform(func() {
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
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.ItemNotFound
		} else {
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
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - m (*T): A pointer to the model to be created.
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
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - m (*T): A pointer to the model to be updated.
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
// Generic Type:
//   - T: The type of the model.
//
// Parameters:
//   - m (*T): A pointer to the model to be deleted.
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
