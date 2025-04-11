package db

import (
	"database/sql"
	"github.com/gflydev/core/utils"
	"github.com/jiveio/fluentsql"
	"github.com/jmoiron/sqlx"
	"log"
)

// ====================================================================
//                              DB Model
// ====================================================================

// Raw struct represents a raw SQL query with its arguments.
type Raw struct {
	sqlStr string // The SQL query string.
	args   []any  // The arguments for the SQL query.
}

// DBModel struct represents a database model with SQL builders and transaction handling.
type DBModel struct {
	tx *sqlx.Tx // Database transaction.

	model any // Model struct for queries.
	raw   Raw // Raw SQL query and arguments.

	selectStatement      fluentsql.Select  // SQL SELECT statement builder.
	omitsSelectStatement fluentsql.Select  // SQL SELECT statement builder for omitting columns.
	whereStatement       fluentsql.Where   // WHERE clause builder.
	joinStatement        fluentsql.Join    // JOIN clause builder.
	groupByStatement     fluentsql.GroupBy // GROUP BY clause builder.
	havingStatement      fluentsql.Having  // HAVING clause builder.
	orderByStatement     fluentsql.OrderBy // ORDER BY clause builder.
	limitStatement       fluentsql.Limit   // LIMIT clause builder.
	fetchStatement       fluentsql.Fetch   // FETCH clause builder, a version of LIMIT.
}

// Instance creates and returns a new DBModel instance.
//
// Returns:
//
//	*DBModel - A pointer to a new database model instance.
func Instance() *DBModel {
	return &DBModel{
		tx:    nil,
		model: nil,
	}
}

// reset clears the state of the DBModel and resets builders.
//
// Returns:
//
//	*DBModel - The reset DBModel instance.
func (db *DBModel) reset() *DBModel {
	db.model = nil                                          // Clear the model.
	db.raw.sqlStr = ""                                      // Reset raw SQL string.
	db.selectStatement.Columns = []any{}                    // Clear SELECT columns.
	db.omitsSelectStatement.Columns = []any{}               // Clear omitted SELECT columns.
	db.whereStatement.Conditions = []fluentsql.Condition{}  // Clear WHERE conditions.
	db.joinStatement.Items = []fluentsql.JoinItem{}         // Clear JOIN items.
	db.groupByStatement.Items = []string{}                  // Clear GROUP BY items.
	db.havingStatement.Conditions = []fluentsql.Condition{} // Clear HAVING conditions.
	db.orderByStatement.Items = []fluentsql.SortItem{}      // Clear ORDER BY items.
	db.limitStatement.Limit = 0                             // Reset limit.
	db.fetchStatement.Fetch = 0                             // Reset fetch.

	return db
}

// ====================================================================
//                      FluentSQL + SQLX integration
// ====================================================================

// get performs fetching a single data row using QueryBuilder.
//
// Parameters:
//   - q (*fluentsql.QueryBuilder): The query builder comprising the SQL query and arguments.
//   - model (any): The model to map the resulting row.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) get(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.getRaw(sqlStr, args, model)
}

// getRaw executes a raw SQL query to fetch a single data row.
//
// Parameters:
//   - sqlStr (string): The raw SQL query string.
//   - args ([]any): Arguments for the query placeholders.
//   - model (any): The model to map the resulting row.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) getRaw(sqlStr string, args []any, model any) (err error) {
	if utils.Getenv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	if db.tx != nil {
		err = db.tx.Get(model, sqlStr, args...)
	} else {
		err = dbInstance.Get(model, sqlStr, args...)
	}

	return
}

// query performs querying a list of data rows using QueryBuilder.
//
// Parameters:
//   - q (*fluentsql.QueryBuilder): The query builder with the SQL and arguments.
//   - model (any): The model to map the resulting rows.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) query(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.queryRaw(sqlStr, args, model)
}

// queryRaw executes a raw SQL query to fetch a list of data rows.
//
// Parameters:
//   - sqlStr (string): The raw SQL query string.
//   - args ([]any): Arguments for the query placeholders.
//   - model (any): The model to map the resulting rows.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) queryRaw(sqlStr string, args []any, model any) (err error) {
	if utils.Getenv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	if db.tx != nil {
		err = db.tx.Select(model, sqlStr, args...)
	} else {
		err = dbInstance.Select(model, sqlStr, args...)
	}

	return
}

// add performs inserting new data using InsertBuilder and returns the inserted ID.
//
// Parameters:
//   - q (*fluentsql.InsertBuilder): The insert query builder with the SQL and arguments.
//   - primaryColumn (string): The primary column to return, used for PostgreSQL.
//
// Returns:
//   - id (any): The ID of the newly inserted row.
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) add(q *fluentsql.InsertBuilder, primaryColumn string) (id any, err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.addRaw(sqlStr, args, primaryColumn)
}

// addRaw executes a raw SQL query to insert new data and returns the inserted ID.
//
// Parameters:
//   - sqlStr (string): The raw SQL insert query string.
//   - args ([]any): Arguments for the query placeholders.
//   - primaryColumn (string): The primary column to return, used for PostgreSQL.
//
// Returns:
//   - id (any): The ID of the newly inserted row.
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) addRaw(sqlStr string, args []any, primaryColumn string) (id any, err error) {
	if utils.Getenv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	// Data persistence
	if fluentsql.DBType() == fluentsql.PostgreSQL {
		if primaryColumn != "" {
			sqlStr += " RETURNING " + primaryColumn

			if utils.Getenv("DB_DEBUG", false) {
				log.Printf("Changed SQL> %s", sqlStr)
			}
		}

		if db.tx != nil {
			err = db.tx.QueryRow(sqlStr, args...).Scan(&id)
		} else {
			err = dbInstance.QueryRow(sqlStr, args...).Scan(&id)
		}
	} else if fluentsql.DBType() == fluentsql.MySQL {
		var result sql.Result
		if db.tx != nil {
			result, _ = db.tx.Exec(sqlStr, args...)
		} else {
			result, _ = dbInstance.Exec(sqlStr, args...)
		}

		id, err = result.LastInsertId()
	}

	return
}

// update performs updating data using UpdateBuilder.
//
// Parameters:
//   - q (*fluentsql.UpdateBuilder): The update query builder with the SQL and arguments.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) update(q *fluentsql.UpdateBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// delete performs deleting data using DeleteBuilder.
//
// Parameters:
//   - q (*fluentsql.DeleteBuilder): The delete query builder with the SQL and arguments.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) delete(q *fluentsql.DeleteBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// execRaw performs executing a raw SQL query for update or delete operations.
//
// Parameters:
//   - sqlStr (string): The raw SQL query string.
//   - args ([]any): Arguments for the query placeholders.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) execRaw(sqlStr string, args []any) (err error) {
	if utils.Getenv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	// Data persistence
	if db.tx != nil {
		_, err = db.tx.Exec(sqlStr, args...)
	} else {
		_, err = dbInstance.Exec(sqlStr, args...)
	}

	return
}

// count retrieves the total number of rows based on the QueryBuilder.
//
// Parameters:
//   - q (*fluentsql.QueryBuilder): The query builder with the SQL and arguments.
//   - total (*int): Pointer to an integer to store the total count.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) count(q *fluentsql.QueryBuilder, total *int) error {
	var fetch fluentsql.Fetch
	var limit fluentsql.Limit

	// Build SQL without pagination
	fetch = q.RemoveFetch()
	limit = q.RemoveLimit()

	// Create COUNT query
	sqlBuilderCount := fluentsql.QueryInstance().
		Select("COUNT(*) AS total").
		From(q, "_result_out_")

	err := db.get(sqlBuilderCount, total)
	if err != nil {
		return err
	}

	// Reset pagination
	q.Limit(limit.Limit, limit.Offset)
	q.Fetch(fetch.Offset, fetch.Fetch)

	return nil
}

// ====================================================================
//                           DB Model operators
// ====================================================================

// Raw builds a query from raw SQL.
//
// Parameters:
//   - sqlStr (string): The raw SQL query string.
//   - args (...any): Variadic arguments for the query placeholders.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Raw(sqlStr string, args ...any) *DBModel {
	db.raw.sqlStr = sqlStr
	db.raw.args = args

	return db
}

// Select specifies the list of columns to retrieve.
//
// Parameters:
//   - columns (...any): Variadic list of columns to include in the SELECT clause.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Select(columns ...any) *DBModel {
	db.selectStatement.Columns = columns

	return db
}

// Omit excludes specific columns from retrieval.
//
// Parameters:
//   - columns (...any): Variadic list of columns to omit.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Omit(columns ...any) *DBModel {
	db.omitsSelectStatement.Columns = columns

	return db
}

// Model sets a specific model for the builder.
//
// Parameters:
//   - model (any): The model instance to operate on.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Model(model any) *DBModel {
	db.model = model

	return db
}

// Where adds a WHERE condition to the query.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (fluentsql.WhereOpt): The operator to use (e.g., equals, greater than).
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Where(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// WhereOr adds an OR condition to the WHERE clause.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (fluentsql.WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) WhereOr(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.Or,
	})

	return db
}

// WhereGroup combines multiple WHERE conditions into a group.
//
// Parameters:
//   - groupCondition (fluentsql.FnWhereBuilder): The function to build grouped conditions.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) WhereGroup(groupCondition fluentsql.FnWhereBuilder) *DBModel {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	cond := fluentsql.Condition{
		Group: whereBuilder.Conditions(),
	}

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, cond)

	return db
}

// When conditionally applies a WHERE condition if the provided condition is TRUE.
//
// Parameters:
//   - condition (bool): Determines whether the condition should be applied.
//   - groupCondition (fluentsql.FnWhereBuilder): The function to build the condition.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) When(condition bool, groupCondition fluentsql.FnWhereBuilder) *DBModel {
	if !condition {
		return db
	}

	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, whereBuilder.Conditions()...)

	return db
}

// Join adds a JOIN clause to the query.
//
// Parameters:
//   - join (fluentsql.JoinType): The type of join (e.g., INNER, LEFT).
//   - table (string): The name of the table to join.
//   - condition (fluentsql.Condition): The condition for the join.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Join(join fluentsql.JoinType, table string, condition fluentsql.Condition) *DBModel {
	db.joinStatement.Append(fluentsql.JoinItem{
		Join:      join,
		Table:     table,
		Condition: condition,
	})

	return db
}

// Having adds a HAVING condition to the query.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (fluentsql.WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Having(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.havingStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// GroupBy adds GROUP BY fields to the query.
//
// Parameters:
//   - fields (...string): The fields to group by.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) GroupBy(fields ...string) *DBModel {
	db.groupByStatement.Append(fields...)

	return db
}

// OrderBy adds an ORDER BY clause to the query.
//
// Parameters:
//   - field (string): The field to sort by.
//   - dir (fluentsql.OrderByDir): The sorting direction (e.g., ASC or DESC).
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) OrderBy(field string, dir fluentsql.OrderByDir) *DBModel {
	db.orderByStatement.Append(field, dir)

	return db
}

// Limit adds a LIMIT clause to the query.
//
// Parameters:
//   - limit (int): The maximum number of rows to return.
//   - offset (int): The number of rows to skip.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Limit(limit, offset int) *DBModel {
	db.limitStatement.Limit = limit
	db.limitStatement.Offset = offset

	return db
}

// RemoveLimit removes the LIMIT clause from the query.
//
// Returns:
//   - fluentsql.Limit: The removed limit settings.
func (db *DBModel) RemoveLimit() fluentsql.Limit {
	var _limitStatement fluentsql.Limit

	_limitStatement.Limit = db.limitStatement.Limit
	_limitStatement.Offset = db.limitStatement.Offset

	db.limitStatement.Limit = 0
	db.limitStatement.Offset = 0

	return _limitStatement
}

// Fetch adds a FETCH clause to the query.
//
// Parameters:
//   - offset (int): The offset for fetching rows.
//   - fetch (int): The number of rows to fetch.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Fetch(offset, fetch int) *DBModel {
	db.fetchStatement.Offset = offset
	db.fetchStatement.Fetch = fetch

	return db
}

// RemoveFetch removes the FETCH clause from the query.
//
// Returns:
//   - fluentsql.Fetch: The removed fetch settings.
func (db *DBModel) RemoveFetch() fluentsql.Fetch {
	var _fetchStatement fluentsql.Fetch

	_fetchStatement.Offset = db.fetchStatement.Offset
	_fetchStatement.Fetch = db.fetchStatement.Fetch

	db.fetchStatement.Offset = 0
	db.fetchStatement.Fetch = 0

	return _fetchStatement
}

// whereFromModel builds and appends a WHERE clause from the specific model's data.
//
// Parameters:
//   - queryBuilder (*fluentsql.QueryBuilder): The query builder to modify.
func (tbl *Table) whereFromModel(queryBuilder *fluentsql.QueryBuilder) {
	if tbl.HasData {
		for _, column := range tbl.Columns {
			// Prevent processing meta, relational, and default (zero) column values
			if column.isNotData() || column.IsZero {
				continue
			}

			// Append query conditions
			queryBuilder.Where(column.Name, fluentsql.Eq, tbl.Values[column.Name])
		}
	}
}

// ====================================================================
//                            DB Transaction
// ====================================================================

// Begin starts a new database transaction.
//
// Returns:
//   - *DBModel: The DBModel instance with an active transaction.
func (db *DBModel) Begin() *DBModel {
	// Initialize a new transaction for the database.
	db.tx = dbInstanceTx()

	return db
}

// Rollback rolls back the current database transaction.
//
// Returns:
//   - error: An error, if any, that occurred during the rollback process.
func (db *DBModel) Rollback() error {
	// Check if there’s an active transaction.
	if db.tx != nil {
		// Attempt to roll back the transaction and return the result.
		return db.tx.Rollback()
	}

	// Return nil if there’s no active transaction.
	return nil
}

// Commit commits the current database transaction.
//
// Returns:
//   - error: An error, if any, that occurred during the commit process.
func (db *DBModel) Commit() error {
	// Check if there’s an active transaction.
	if db.tx != nil {
		// Attempt to commit the transaction and return the result.
		return db.tx.Commit()
	}

	// Return nil if there’s no active transaction.
	return nil
}
