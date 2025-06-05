package db

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	qb "github.com/jivegroup/fluentsql"
	"reflect"
	"slices"
)

// Create adds new data for a table via model type Slice, Struct, or *Struct.
//
// Parameters:
//   - model: The data to be inserted. Accepts Slice, Map, Struct, or *Struct types.
//
// Returns:
//   - error: Returns an error if the operation fails.
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
		log.Fatal(err)
	}

	// Reset fluent model builder
	db.reset()

	return
}

// createByRaw inserts data using a raw SQL query string.
//
// Parameters:
//   - model: The data model to insert.
//
// Returns:
//   - error: Returns an error if the operation fails.
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

// createByMap inserts data by reflecting keys and values from a Map.
//
// Parameters:
//   - value: The Map containing the keys and values to insert.
//
// Returns:
//   - error: Returns an error if the operation fails.
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

// createBySlice inserts multiple records by reflecting over a Slice of models.
//
// Parameters:
//   - model: The Slice containing the models to insert.
//
// Returns:
//   - error: Returns an error if the operation fails.
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

// createByStruct inserts a single record using reflection on a Struct.
//
// Parameters:
//   - model: The Struct model to insert.
//
// Returns:
//   - error: Returns an error if the operation fails.
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
