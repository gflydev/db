package db

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/jiveio/fluentsql"
	"reflect"
)

// Update modify data for table via model type Struct, *Struct
func (db *DBModel) Update(model any) (err error) {
	typ := reflect.TypeOf(model)

	switch {
	case db.raw.sqlStr != "":
		err = db.execRaw(db.raw.sqlStr, db.raw.args)
	case typ.Kind() == reflect.Map:
		err = db.updateByMap(model)
	case typ.Kind() == reflect.Struct ||
		(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct):
		err = db.updateByStruct(model)
	}

	if err != nil {
		panic(err)
	}

	// Reset fluent model builder
	db.reset()

	return
}

func (db *DBModel) updateByMap(value any) error {
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

	err = db.updateByStruct(db.model)

	return err
}

// updateByStruct Update data by reflection Struct
func (db *DBModel) updateByStruct(model any) (err error) {
	var table *Table
	var hasCondition = false

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		panic(err)
	}

	// Create update instance
	updateBuilder := fluentsql.UpdateInstance().
		Update(table.Name)

	// Build WHERE condition from a condition list
	for _, condition := range db.whereStatement.Conditions {
		// Sub-conditions
		switch {
		case len(condition.Group) > 0:
			// Append conditions from a group to query builder
			updateBuilder.WhereGroup(func(whereBuilder fluentsql.WhereBuilder) *fluentsql.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)

				return &whereBuilder
			})
			hasCondition = true
		case condition.AndOr == fluentsql.And:
			// Add Where AND condition
			updateBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == fluentsql.Or:
			// Add Where OR condition
			updateBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		}
	}

	// Build WHERE condition with primary column values
	if !hasCondition {
		for _, column := range table.Columns {
			if column.Primary {
				value := table.Values[column.Name]

				updateBuilder.Where(column.Name, fluentsql.Eq, value)
				hasCondition = true
			}
		}
	}

	if !hasCondition {
		panic(errors.New("missing WHERE condition for updating operator"))
	}

	// Build Updating fields from model's data
	for _, column := range table.Columns {
		if column.isNotData() || column.Primary {
			continue
		}

		// Append SET values
		updateBuilder.Set(column.Name, table.Values[column.Name])
	}

	err = db.update(updateBuilder)

	return
}
