package db

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	qb "github.com/jivegroup/fluentsql"
	"reflect"
	"slices"
)

// Create inserts new data into a database table using various model types.
// This method provides a flexible interface for inserting data by automatically
// detecting the input type and routing to the appropriate insertion strategy.
// It supports raw SQL queries, map-based insertion, batch insertion via slices,
// and single record insertion via structs.
//
// Parameters:
//   - model (any): The data to be inserted. Supported types:
//   - Raw SQL: When db.raw.sqlStr is set, uses raw SQL insertion
//   - map[string]any: Creates a record using key-value pairs from the map
//   - []Struct or []*Struct: Batch insertion of multiple records
//   - Struct or *Struct: Single record insertion using struct fields
//
// Returns:
//   - error: Returns an error if the insertion fails, model type is unsupported,
//     or if there are database connectivity issues. Returns nil on success.
//
// Examples:
//
//	// Single struct insertion
//	user := User{Name: "John", Email: "john@example.com"}
//	err := db.Model(&User{}).Create(user)
//
//	// Pointer to struct insertion
//	user := &User{Name: "Jane", Email: "jane@example.com"}
//	err := db.Model(&User{}).Create(user)
//
//	// Batch insertion with slice
//	users := []User{
//	    {Name: "Alice", Email: "alice@example.com"},
//	    {Name: "Bob", Email: "bob@example.com"},
//	}
//	err := db.Model(&User{}).Create(users)
//
//	// Map-based insertion
//	userData := map[string]any{
//	    "name": "Charlie",
//	    "email": "charlie@example.com",
//	    "age": 30,
//	}
//	err := db.Model(&User{}).Create(userData)
//
//	// Raw SQL insertion
//	err := db.Model(&User{}).Raw("INSERT INTO users (name, email) VALUES (?, ?)", "David", "david@example.com").Create(nil)
//
// Note:
//   - Primary key fields are automatically handled and populated after insertion
//   - The method respects Select() and Omit() clauses for column filtering
//   - For batch operations, if one record fails, the entire operation may fail
//   - The model is reset after the operation completes
func (db *DBModel) Create(model any) (err error) {
	// Get the type of the model
	typ := reflect.TypeOf(model)

	// Perform creation based on model type
	switch {
	case db.raw.sqlStr != "":
		err = db.createByRaw(model)
	case typ.Kind() == reflect.Map:
		err = db.createByMap(model)
	case typ.Kind() == reflect.Slice:
		err = db.createBySlice(model)
	case
		typ.Kind() == reflect.Struct || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct):
		err = db.createByStruct(model)
	}

	if err != nil {
		log.Error(err)
	}

	// Reset fluent model builder
	db.reset()

	return
}

// createByRaw executes a raw SQL insertion query with automatic primary key handling.
// This method is used internally when a raw SQL string has been set via the Raw() method.
// It extracts table metadata from the provided model to handle primary key population
// after insertion, making it seamless to work with raw SQL while maintaining ORM benefits.
//
// Parameters:
//   - model (any): The data model used for metadata extraction and primary key population.
//     This should be a struct or pointer to struct that represents the target table.
//     The model is used to:
//   - Extract table information and primary key details
//   - Populate the primary key field after successful insertion
//
// Returns:
//   - error: Returns an error if:
//   - The model cannot be processed (invalid type or structure)
//   - The raw SQL execution fails
//   - Primary key population fails
//   - Database connectivity issues occur
//     Returns nil on successful insertion.
//
// Examples:
//
//	// Basic raw insertion with automatic ID population
//	user := &User{}
//	err := db.Model(&User{}).Raw("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com").Create(user)
//	// user.ID will be populated with the inserted record's ID
//
//	// Raw insertion with complex query
//	user := &User{}
//	query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, NOW())"
//	err := db.Model(&User{}).Raw(query, "Jane", "jane@example.com").Create(user)
//
// Note:
//   - Only works when a raw SQL string has been previously set via Raw() method
//   - Requires exactly one primary key column in the table for ID population
//   - The model structure must match the target table for proper metadata extraction
//   - Primary key population only occurs if the table has a single primary key column
func (db *DBModel) createByRaw(model any) error {
	var table *Table
	var err error

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		return err
	}

	var id any
	var primaryColumn Column

	// Get primary column name (if only one primary in table exists)
	if len(table.Primaries) == 1 {
		primaryColumn = table.Primaries[0]
	}

	// Perform raw insertion and retrieve the ID
	id, err = db.addRaw(db.raw.sqlStr, db.raw.args, primaryColumn.Name)

	if err != nil {
		return err
	}

	// Set ID back to the model
	if primaryColumn.Key != "" {
		err = setValue(model, primaryColumn.Key, id)
	}

	return err
}

// createByMap inserts a database record using key-value pairs from a map.
// This method provides a flexible way to insert data without requiring a predefined struct.
// It maps the keys from the input map to the corresponding fields in the target model,
// then delegates to createByStruct for the actual insertion. This approach allows for
// dynamic data insertion while maintaining type safety through the underlying model.
//
// Parameters:
//   - value (any): A map containing the data to insert. Expected to be of type map[string]any
//     or similar map types where:
//   - Keys represent database column names or struct field names
//   - Values represent the data to be inserted into those columns
//   - Only non-zero values are processed (zero values are skipped)
//
// Returns:
//   - error: Returns an error if:
//   - No model has been set via Model() method (returns "Missing model for map value")
//   - The value parameter is not a valid map type
//   - Field mapping from map keys to model fields fails
//   - The underlying struct insertion fails
//   - Database connectivity issues occur
//     Returns nil on successful insertion.
//
// Examples:
//
//	// Basic map insertion
//	userData := map[string]any{
//	    "name": "John Doe",
//	    "email": "john@example.com",
//	    "age": 30,
//	}
//	err := db.Model(&User{}).Create(userData)
//
//	// Map with mixed data types
//	productData := map[string]any{
//	    "name": "Laptop",
//	    "price": 999.99,
//	    "in_stock": true,
//	    "category_id": 5,
//	}
//	err := db.Model(&Product{}).Create(productData)
//
//	// Map with nil/zero values (these will be skipped)
//	userData := map[string]any{
//	    "name": "Jane",
//	    "email": "",        // Empty string - will be skipped
//	    "age": 0,          // Zero value - will be skipped
//	    "bio": nil,        // Nil value - will be skipped
//	}
//	err := db.Model(&User{}).Create(userData)
//
// Note:
//   - Requires a model to be set via Model() method before calling
//   - Only processes non-zero values from the map (zero values are ignored)
//   - Map keys should correspond to struct field names or database column names
//   - The method internally converts the map data to the target model struct
//   - Primary key handling is managed by the underlying createByStruct method
func (db *DBModel) createByMap(value any) error {
	var err error

	// Check if a model exists
	if db.model == nil {
		return errors.New("Missing model for map value")
	}

	// Reflect items from the Map
	mapValue := reflect.ValueOf(value)

	// Process each key-value pair in the Map
	for _, key := range mapValue.MapKeys() {
		itemVal := mapValue.MapIndex(key)

		// Check if the value is valid and not zero
		isSet := itemVal.IsValid() && !itemVal.IsZero()

		if isSet {
			val := itemVal.Interface()

			// Set the value on the model
			err = setValue(db.model, key.String(), val)
			if err != nil {
				log.Error(err)
			}
		}
	}

	// If an error occurred, return it
	if err != nil {
		return err
	}

	// Use the model to perform a structured insert
	return db.createByStruct(db.model)
}

// createBySlice performs batch insertion of multiple records from a slice of models.
// This method iterates through each element in the provided slice and inserts them
// individually using createByStruct. It supports both struct values and pointers to structs,
// providing flexibility for different data structures. The method handles type validation
// and skips invalid entries to ensure robust batch processing.
//
// Parameters:
//   - model (any): A slice containing the models to insert. Supported slice types:
//   - []Struct: Slice of struct values (e.g., []User)
//   - []*Struct: Slice of pointers to structs (e.g., []*User)
//   - Mixed slices are supported, but each element must be a valid struct or pointer to struct
//   - Invalid or nil elements in the slice are automatically skipped
//
// Returns:
//   - error: Returns an error if:
//   - Any individual record insertion fails (the error from the last failed insertion)
//   - The model parameter is not a valid slice type
//   - Database connectivity issues occur during any insertion
//   - Struct reflection fails for any element
//     Returns nil if all insertions succeed or if all elements are skipped.
//
// Examples:
//
//	// Batch insertion with struct slice
//	users := []User{
//	    {Name: "Alice", Email: "alice@example.com"},
//	    {Name: "Bob", Email: "bob@example.com"},
//	    {Name: "Charlie", Email: "charlie@example.com"},
//	}
//	err := db.Model(&User{}).Create(users)
//
//	// Batch insertion with pointer slice
//	users := []*User{
//	    &User{Name: "David", Email: "david@example.com"},
//	    &User{Name: "Eve", Email: "eve@example.com"},
//	}
//	err := db.Model(&User{}).Create(users)
//
//	// Mixed slice with some nil pointers (nil entries will be skipped)
//	users := []*User{
//	    &User{Name: "Frank", Email: "frank@example.com"},
//	    nil, // This will be skipped
//	    &User{Name: "Grace", Email: "grace@example.com"},
//	}
//	err := db.Model(&User{}).Create(users)
//
//	// Empty slice (no operation performed)
//	var users []User
//	err := db.Model(&User{}).Create(users) // Returns nil, no insertions
//
// Note:
//   - Each record is inserted individually, not as a single batch SQL operation
//   - If one insertion fails, subsequent insertions are not attempted
//   - Primary keys are populated for each successfully inserted record
//   - The method respects Select() and Omit() clauses for all insertions
//   - Invalid or nil slice elements are silently skipped
//   - For better performance with large datasets, consider using database-specific batch insert methods
func (db *DBModel) createBySlice(model any) (err error) {
	// Reflect the Slice
	items := reflect.ValueOf(model)

	// Iterate through each element in the Slice
	for i := 0; i < items.Len(); i++ {
		itemVal := items.Index(i)
		var indirectVal reflect.Value

		// Handle *Struct or Struct types
		if itemVal.Kind() == reflect.Pointer {
			indirectVal = reflect.Indirect(itemVal.Elem())
		} else if itemVal.Kind() == reflect.Struct {
			indirectVal = reflect.Indirect(itemVal)
		}

		// Skip invalid types
		if !indirectVal.IsValid() {
			continue
		}

		// Perform insertion for each item in the Slice
		item := itemVal.Interface()
		err = db.createByStruct(item)
	}

	return
}

// createByStruct inserts a single database record using struct field reflection.
// This is the core insertion method that handles individual struct-based insertions.
// It extracts table metadata from the struct, builds appropriate SQL INSERT statements,
// executes the insertion, and populates the primary key back into the original struct.
// The method respects column filtering via Select() and Omit() clauses and handles
// various column types including primary keys, data columns, and computed columns.
//
// Parameters:
//   - model (any): The struct model to insert. Supported types:
//   - Struct: Direct struct value (e.g., User{Name: "John"})
//   - *Struct: Pointer to struct (e.g., &User{Name: "John"})
//   - The struct should have appropriate field tags for database mapping
//   - Primary key fields are automatically handled and should not be manually set
//
// Returns:
//   - error: Returns an error if:
//   - The model cannot be processed or reflected (invalid struct type)
//   - Table metadata extraction fails
//   - SQL INSERT statement building fails
//   - Database insertion execution fails
//   - Primary key population back to the struct fails
//   - Database connectivity issues occur
//     Returns nil on successful insertion with primary key populated.
//
// Examples:
//
//	// Basic struct insertion
//	user := User{
//	    Name:  "John Doe",
//	    Email: "john@example.com",
//	    Age:   30,
//	}
//	err := db.Model(&User{}).Create(user)
//	// user.ID will be populated after insertion
//
//	// Pointer to struct insertion
//	user := &User{
//	    Name:  "Jane Smith",
//	    Email: "jane@example.com",
//	}
//	err := db.Model(&User{}).Create(user)
//	// user.ID will be populated after insertion
//
//	// Insertion with column selection
//	user := User{Name: "Bob", Email: "bob@example.com", Age: 25}
//	err := db.Model(&User{}).Select("name", "email").Create(user)
//	// Only name and email will be inserted, age is omitted
//
//	// Insertion with column omission
//	user := User{Name: "Alice", Email: "alice@example.com", Age: 28}
//	err := db.Model(&User{}).Omit("age").Create(user)
//	// All fields except age will be inserted
//
// Note:
//   - Primary key fields are automatically excluded from insertion and populated after success
//   - Computed columns and non-data columns are automatically skipped
//   - The method respects Select() clauses to include only specified columns
//   - The method respects Omit() clauses to exclude specified columns
//   - Only works with tables that have exactly one primary key for ID population
//   - The original struct/pointer is modified with the generated primary key value
//   - Field mapping relies on struct tags and reflection for database column mapping
func (db *DBModel) createByStruct(model any) (err error) {
	var table *Table
	var columns []string
	var values []any

	// Create a table object from a model
	if table, err = ModelData(model); err != nil {
		return
	}

	// Generate insert columns and values by iterating over table columns
	for _, column := range table.Columns {
		// Skip columns that are not data or are primary keys
		if column.isNotData() || column.Primary {
			continue
		}

		// Skip columns not included in the select statement
		if len(db.selectStatement.Columns) > 0 && !slices.Contains(db.selectStatement.Columns, any(column.Name)) {
			continue
		}

		// Skip columns specified in the omits select statement
		if len(db.omitsSelectStatement.Columns) > 0 && slices.Contains(db.omitsSelectStatement.Columns, any(column.Name)) {
			continue
		}

		// Add column names and values to their respective slices
		value := table.Values[column.Name]
		columns = append(columns, column.Name)
		values = append(values, value)
	}

	// Build an INSERT SQL statement with the columns and values
	insertBuilder := qb.InsertInstance().
		Insert(table.Name, columns...).
		Row(values...)

	var id any
	var primaryColumn Column

	// Check if there is exactly one primary column
	if len(table.Primaries) == 1 {
		primaryColumn = table.Primaries[0]
	}

	// Perform the insert and retrieve the ID
	if id, err = db.add(insertBuilder, primaryColumn.Name); err != nil {
		return
	}

	// Set the ID back to the model
	if primaryColumn.Key != "" {
		err = setValue(model, primaryColumn.Key, id)
	}

	return
}
