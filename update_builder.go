package db

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	qb "github.com/jivegroup/fluentsql"
	"reflect"
)

// Update modifies data for a table using a model of type Struct or *Struct.
//
// Parameters:
//   - model (any): The data model used for updating the table. It can be of type Struct, *Struct, or a map.
//
// Returns:
//   - error: Returns an error if the update process fails.
func (db *DBModel) Update(model any) (err error) {
	typ := reflect.TypeOf(model)

	switch {
	case db.raw.sqlStr != "":
		// Execute raw SQL query if provided
		err = db.execRaw(db.raw.sqlStr, db.raw.args)
	case typ.Kind() == reflect.Map:
		// Update using map data
		err = db.updateByMap(model)
	case typ.Kind() == reflect.Struct || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct):
		// Update using struct data
		err = db.updateByStruct(model)
	}

	// Reset fluent model builder
	db.reset()

	return
}

// updateByMap updates data in the database when the provided model is of type map.
//
// Parameters:
//   - value (any): A map where keys correspond to field names and values to the data to be updated.
//
// Returns:
//   - error: Returns an error if the update process fails.
func (db *DBModel) updateByMap(value any) error {
	var err error

	if db.model == nil {
		return errors.New("Missing model for map value")
	}

	// Reflect items from the map
	mapValue := reflect.ValueOf(value)

	// Process each map key and update corresponding model fields
	for _, key := range mapValue.MapKeys() {
		itemVal := mapValue.MapIndex(key)

		// Check if the value is valid and not zero
		isSet := itemVal.IsValid() && !itemVal.IsZero()

		if isSet {
			// Convert map value to interface
			val := itemVal.Interface()

			// Set the value on the model using reflection
			err = setValue(db.model, key.String(), val)
			if err != nil {
				// Log the error
				log.Error(err)
			}
		}
	}

	// Delegate the updated model to the struct update function
	return db.updateByStruct(db.model)
}

// updateByStruct Update modifies table data using a data model.
//
// Parameters:
//   - model (any): The data model, which can be a Struct or *Struct.
//
// Returns:
//   - error: Returns an error if the update process fails or if a WHERE condition is missing.
func (db *DBModel) updateByStruct(model any) (err error) {
	var (
		table        *Table // Table representation created from the model.
		hasCondition bool   // Flag indicating whether a WHERE condition exists.
	)

	// Create a table object from the data model.
	if table, err = ModelData(model); err != nil {
		return
	}

	// Initialize the Update query builder for the target database table.
	updateBuilder := qb.UpdateInstance().Update(table.Name)

	// Build WHERE conditions from pre-defined conditions in 'whereStatement'.
	for _, condition := range db.whereStatement.Conditions {
		switch {
		case len(condition.Group) > 0:
			// Append grouped conditions using a builder function.
			updateBuilder.WhereGroup(func(whereBuilder qb.WhereBuilder) *qb.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)
				return &whereBuilder
			})
			hasCondition = true
		case condition.AndOr == And:
			// Add an AND clause to the WHERE condition.
			updateBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == Or:
			// Add an OR clause to the WHERE condition.
			updateBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		}
	}

	// Build WHERE condition using primary key column values if no other condition exists.
	if !hasCondition {
		for _, column := range table.Columns {
			if column.Primary {
				// Use primary key column value for the WHERE condition.
				value := table.Values[column.Name]

				updateBuilder.Where(column.Name, Eq, value)
				hasCondition = true
			}
		}
	}

	// Panic if no WHERE condition exists to prevent accidental updates to all rows.
	if !hasCondition {
		err = errors.New("missing WHERE condition for updating operator")
		return
	}

	// Iterate through the table's columns and add SET clauses for valid data fields.
	for _, column := range table.Columns {
		// Skip processing for columns that are not valid data fields or are primary keys.
		if column.isNotData() || column.Primary {
			continue
		}

		// Append a SET clause with the column name and its corresponding value.
		updateBuilder.Set(column.Name, table.Values[column.Name])
	}

	// Execute the update operation using the constructed query builder.
	err = db.update(updateBuilder)

	return
}
