package db

import (
	"github.com/gflydev/core/errors"
	"github.com/jiveio/fluentsql"
)

// Delete perform delete data for table via model type Struct, *Struct
func (db *DBModel) Delete(model any) error {
	var err error

	// Delete by raw sql
	if db.raw.sqlStr != "" {
		err = db.execRaw(db.raw.sqlStr, db.raw.args)

		if err != nil {
			panic(err)
		}

		// Reset fluent model builder
		db.reset()
	}

	var table *Table
	var hasCondition = false

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		panic(err)
	}

	// Create delete instance
	deleteBuilder := fluentsql.DeleteInstance().
		Delete(table.Name)

	// Build WHERE condition with primary columns
	for _, primaryColumn := range table.Primaries {
		primaryKey := primaryColumn.Name
		primaryVal := table.Values[primaryKey]

		if primaryVal != nil {
			wherePrimaryCondition := fluentsql.Condition{
				Field: primaryKey,
				Opt:   fluentsql.Eq,
				Value: primaryVal,
				AndOr: fluentsql.And,
			}

			deleteBuilder.WhereCondition(wherePrimaryCondition)
			hasCondition = true
		}
	}

	// Build WHERE condition from a condition list
	for _, condition := range db.whereStatement.Conditions {
		switch {
		case len(condition.Group) > 0:
			// Append conditions from a group to query builder
			deleteBuilder.WhereGroup(func(whereBuilder fluentsql.WhereBuilder) *fluentsql.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)

				return &whereBuilder
			})
			hasCondition = true
		case condition.AndOr == fluentsql.And:
			// Add Where AND condition
			deleteBuilder.Where(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		case condition.AndOr == fluentsql.Or:
			// Add Where OR condition
			deleteBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
			hasCondition = true
		}
	}

	if !hasCondition {
		panic(errors.New("Missing WHERE condition for deleting operator"))
	}

	err = db.delete(deleteBuilder)

	// Reset fluent model builder
	db.reset()

	return err
}
