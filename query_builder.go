package db

import (
	"crypto/rand"
	"fmt"
	"github.com/gflydev/core/errors"
	"github.com/jiveio/fluentsql"
	"math/big"
	"reflect"
)

// ====================================================================
//                          Query ONE row
// ====================================================================

// GetOne represents a strategy for retrieving a single record from the database.
// Possible values are GetFirst, GetLast, and TakeOne.
type GetOne int

const (
	// GetFirst retrieves the first record ordered by primary key in ascending order.
	GetFirst GetOne = iota

	// GetLast retrieves the last record ordered by primary key in descending order.
	GetLast

	// TakeOne retrieves a random record.
	TakeOne
)

// First retrieves the first record ordered by primary key.
//
// Parameters:
//   - model (any): A pointer to the model where the result will be stored.
//
// Returns:
//   - err (error): An error object if any issues occur during the retrieval process; nil otherwise.
func (db *DBModel) First(model any) (err error) {
	err = db.Get(model, GetFirst)
	return
}

// Last retrieves the last record ordered by primary key in descending order.
//
// Parameters:
//   - model (any): A pointer to the model where the result will be stored.
//
// Returns:
//   - err (error): An error object if any issues occur during the retrieval process; nil otherwise.
func (db *DBModel) Last(model any) (err error) {
	err = db.Get(model, GetLast)
	return
}

// Get retrieves a single record from the database based on the specified strategy.
//
// Parameters:
//   - model (any): A pointer to the model where the retrieved record will be stored.
//   - getType (GetOne): The strategy for selecting the record. Possible values are:
//   - GetFirst: Retrieve the first record ordered by primary key in ascending order.
//   - GetLast: Retrieve the last record ordered by primary key in descending order.
//   - TakeOne: Retrieve a random record.
//
// Returns:
//   - err (error): An error object if any issues occur during the retrieval process; nil otherwise.
func (db *DBModel) Get(model any, getType GetOne) (err error) {
	// Query raw SQL
	if db.raw.sqlStr != "" {
		// Data persistence
		if db.tx != nil {
			err = db.tx.Get(model, db.raw.sqlStr, db.raw.args...)
		} else {
			err = dbInstance.Get(model, db.raw.sqlStr, db.raw.args...)
		}

		if err != nil {
			panic(err)
		}

		// Reset fluent model builder
		db.reset()

		return
	}

	// Verify the type of the input model
	typ := reflect.TypeOf(model)
	if !(typ.Kind() == reflect.Struct ||
		(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct)) {
		err = errors.New("Invalid data :: model not Struct type")
		return
	}

	var table *Table

	// Create a table object from a model
	table, err = ModelData(model)
	if err != nil {
		panic(err)
	}

	// Define the columns to be queried
	var selectColumns []any
	if len(db.selectStatement.Columns) > 0 {
		selectColumns = db.selectStatement.Columns
	} else {
		selectColumns = []any{"*"}
	}

	// Create a query builder with initial SELECT, FROM, and LIMIT clauses
	queryBuilder := fluentsql.QueryInstance().
		Select(selectColumns...).
		From(table.Name).
		Limit(1, 0)

	// Build WHERE condition using primary columns
	for _, primaryColumn := range table.Primaries {
		primaryKey := primaryColumn.Name
		primaryVal := table.Values[primaryKey]

		if primaryVal != nil {
			// Build WHERE condition with specific primary value
			wherePrimaryCondition := fluentsql.Condition{
				Field: primaryKey,
				Opt:   fluentsql.Eq,
				Value: primaryVal,
				AndOr: fluentsql.And,
			}
			queryBuilder.WhereCondition(wherePrimaryCondition)
		}
	}

	// Build WHERE condition from condition list
	for _, condition := range db.whereStatement.Conditions {
		// Handle grouped or individual conditions
		switch {
		case len(condition.Group) > 0:
			queryBuilder.WhereGroup(func(whereBuilder fluentsql.WhereBuilder) *fluentsql.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)
				return &whereBuilder
			})
		case condition.AndOr == fluentsql.And:
			queryBuilder.Where(condition.Field, condition.Opt, condition.Value)
		case condition.AndOr == fluentsql.Or:
			queryBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
		}
	}

	// Build WHERE condition from the model's data
	table.whereFromModel(queryBuilder)

	// Build JOIN clause
	for _, joinItem := range db.joinStatement.Items {
		queryBuilder.Join(joinItem.Join, joinItem.Table, joinItem.Condition)
	}

	// Build GROUP BY clause
	if len(db.groupByStatement.Items) > 0 {
		queryBuilder.GroupBy(db.groupByStatement.Items...)
	}

	// Build HAVING clause
	for _, condition := range db.havingStatement.Conditions {
		queryBuilder.Having(condition.Field, condition.Opt, condition.Value)
	}

	// Build LIMIT clause
	if db.limitStatement.Limit > 0 {
		queryBuilder.Limit(db.limitStatement.Limit, db.limitStatement.Offset)
	}

	// Build FETCH clause
	if db.fetchStatement.Fetch > 0 {
		queryBuilder.Fetch(db.fetchStatement.Offset, db.fetchStatement.Fetch)
	}

	// Build ORDER BY clause
	orderByField := ""
	if len(table.Primaries) > 0 {
		orderByField = table.Primaries[0].Name
	} else {
		orderByField = table.Columns[0].Name
	}

	var orderByDir fluentsql.OrderByDir

	// Determine order by direction based on the strategy
	switch {
	case getType == GetLast && orderByField != "":
		orderByDir = fluentsql.Desc
	case getType == GetFirst && orderByField != "":
		orderByDir = fluentsql.Asc
	case getType == TakeOne: // Random order by field and direction
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(table.Columns)-1)))
		orderByField = table.Columns[n.Int64()].Name

		n, _ = rand.Int(rand.Reader, big.NewInt(10))
		if n.Int64()%2 == 1 {
			orderByDir = fluentsql.Asc
		} else {
			orderByDir = fluentsql.Desc
		}
	}
	queryBuilder.OrderBy(orderByField, orderByDir)

	// Data processing using the constructed query
	if err = db.get(queryBuilder, model); err != nil {
		panic(err)
	}

	// Reset fluent model builder
	db.reset()

	return
}

// ====================================================================
//                           Query MULTI rows
// ====================================================================

// Find searches for multiple rows in the database based on query criteria.
//
// Parameters:
//   - model (any): A pointer to the slice where the retrieved rows will be stored.
//
// Returns:
//   - total (int): The total number of rows matching the query criteria.
//   - err (error): An error object if any issues occur during the retrieval process; nil otherwise.
func (db *DBModel) Find(model any) (total int, err error) {
	// Query raw SQL
	if db.raw.sqlStr != "" {
		// Data persistence
		if db.tx != nil {
			err = db.tx.Select(model, db.raw.sqlStr, db.raw.args...)
		} else {
			err = dbInstance.Select(model, db.raw.sqlStr, db.raw.args...)
		}

		if err != nil {
			panic(err)
		}

		// Query COUNT
		sqlCount := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s) _result_out_", db.raw.sqlStr)
		err = db.getRaw(sqlCount, db.raw.args, &total)

		if err != nil {
			panic(err)
		}

		// Reset fluent model builder
		db.reset()

		return
	}

	// Validate input type
	typ := reflect.TypeOf(model)
	if !(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice) {
		panic(errors.New("Invalid data :: model not *Slice type"))
	}

	var table *Table

	// Get the type of model and create a table representation
	typeElement := reflect.TypeOf(model).Elem().Elem()  // First Elem() for pointer, second Elem() for item
	valueElement := reflect.ValueOf(typeElement).Elem() // Create empty value
	table = processModel(typeElement, valueElement, NewTable())

	// Define the columns to query
	var selectColumns []any
	if len(db.selectStatement.Columns) > 0 {
		selectColumns = db.selectStatement.Columns
	} else {
		selectColumns = []any{"*"}
	}

	// Create query builder
	queryBuilder := fluentsql.QueryInstance().
		Select(selectColumns...).
		From(table.Name)

	// Build WHERE condition from the condition list
	for _, condition := range db.whereStatement.Conditions {
		// Sub-conditions
		switch {
		case len(condition.Group) > 0:
			// Append conditions from a group to query builder
			queryBuilder.WhereGroup(func(whereBuilder fluentsql.WhereBuilder) *fluentsql.WhereBuilder {
				whereBuilder.WhereCondition(condition.Group...)
				return &whereBuilder
			})
		case condition.AndOr == fluentsql.And:
			// Add Where AND condition
			queryBuilder.Where(condition.Field, condition.Opt, condition.Value)
		case condition.AndOr == fluentsql.Or:
			// Add Where OR condition
			queryBuilder.WhereOr(condition.Field, condition.Opt, condition.Value)
		}
	}

	// Build WHERE condition from model's data in the table
	table.whereFromModel(queryBuilder)

	// Build JOIN clause
	for _, joinItem := range db.joinStatement.Items {
		queryBuilder.Join(joinItem.Join, joinItem.Table, joinItem.Condition)
	}

	// Build GROUP BY clause
	if len(db.groupByStatement.Items) > 0 {
		queryBuilder.GroupBy(db.groupByStatement.Items...)
	}

	// Build HAVING clause
	for _, condition := range db.havingStatement.Conditions {
		queryBuilder.Having(condition.Field, condition.Opt, condition.Value)
	}

	// Build LIMIT clause
	if db.limitStatement.Limit > 0 {
		queryBuilder.Limit(db.limitStatement.Limit, db.limitStatement.Offset)
	}

	// Build FETCH clause
	if db.fetchStatement.Fetch > 0 {
		queryBuilder.Fetch(db.fetchStatement.Offset, db.fetchStatement.Fetch)
	}

	// Build ORDER BY clause
	for _, orderItem := range db.orderByStatement.Items {
		queryBuilder.OrderBy(orderItem.Field, orderItem.Direction)
	}

	// Execute query and populate model
	if err = db.query(queryBuilder, model); err != nil {
		panic(err)
	}

	// Execute count query to get the total number of rows
	if err = db.count(queryBuilder, &total); err != nil {
		panic(err)
	}

	// Reset fluent model builder
	db.reset()

	return
}
