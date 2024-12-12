package db

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/jiveio/fluentsql"
	"reflect"
	"slices"
)

// Create add new data for table via model type Slice, Struct, *Struct
func (db *DBModel) Create(model any) (err error) {
	typ := reflect.TypeOf(model)

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

func (db *DBModel) createByRaw(model any) error {
	var table *Table

	// Create a table object from a model
	table, err := ModelData(model)
	if err != nil {
		panic(err)
	}

	var id any
	var primaryColumn Column

	// Get primary column name (in case only one primary in table)
	if len(table.Primaries) == 1 {
		primaryColumn = table.Primaries[0]
	}

	id, err = db.addRaw(db.raw.sqlStr, db.raw.args, primaryColumn.Name)

	// Set ID back model
	if primaryColumn.Key != "" {
		err = setValue(model, primaryColumn.Key, id)
		if err != nil {
			panic(err)
		}
	}

	return err
}

func (db *DBModel) createByMap(value any) error {
	var err error

	if db.model == nil {
		return errors.New("Missing model for map value")
	}

	// Reflect items from a map
	mapValue := reflect.ValueOf(value)

	// Process for each map key
	for _, key := range mapValue.MapKeys() {
		itemVal := mapValue.MapIndex(key)

		// IsZero panics if the value is invalid.
		// Most functions and methods never return an invalid Value.
		isSet := itemVal.IsValid() && !itemVal.IsZero()

		if isSet {
			val := itemVal.Interface()

			err = setValue(db.model, key.String(), val)
			if err != nil {
				log.Error(err)
			}
		}
	}

	if err != nil {
		return err
	}

	err = db.createByStruct(db.model)

	return err
}

// createBySlice Insert data by reflection Slice
func (db *DBModel) createBySlice(model any) (err error) {
	items := reflect.ValueOf(model)

	for i := 0; i < items.Len(); i++ {
		itemVal := items.Index(i)
		var indirectVal reflect.Value

		// Only process 2 types: *Struct or Struct
		if itemVal.Kind() == reflect.Pointer {
			indirectVal = reflect.Indirect(itemVal.Elem())
		} else if itemVal.Kind() == reflect.Struct {
			indirectVal = reflect.Indirect(itemVal)
		}

		// Skip all unknown types (Not *Struct or Struct)
		if !indirectVal.IsValid() {
			continue
		}

		// Convert reflection item value to an Interface type
		item := itemVal.Interface()
		err = db.createByStruct(item)
	}

	return
}

// createByStruct Insert data by reflection Struct
func (db *DBModel) createByStruct(model any) (err error) {
	var table *Table
	var columns []string
	var values []any

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		panic(err)
	}

	// Get columns and values accordingly
	for _, column := range table.Columns {
		// Restriction from model declaration
		if column.isNotData() || column.Primary {
			continue
		}

		// Restriction from selected columns
		if len(db.selectStatement.Columns) > 0 && !slices.Contains(db.selectStatement.Columns, any(column.Name)) {
			continue
		}

		// Restriction from omitted columns
		if len(db.omitsSelectStatement.Columns) > 0 && slices.Contains(db.omitsSelectStatement.Columns, any(column.Name)) {
			continue
		}

		value := table.Values[column.Name]

		// Pair columns and values to insert
		columns = append(columns, column.Name)
		values = append(values, value)
	}

	// Create insert instance
	insertBuilder := fluentsql.InsertInstance().
		Insert(table.Name, columns...).
		Row(values...)

	var id any
	var primaryColumn Column

	// Get primary column name (in case only one primary in table)
	if len(table.Primaries) == 1 {
		primaryColumn = table.Primaries[0]
	}

	id, err = db.add(insertBuilder, primaryColumn.Name)
	if err != nil {
		panic(err)
	}

	// Set ID back model
	if primaryColumn.Key != "" {
		err = setValue(model, primaryColumn.Key, id)
	}

	return
}
