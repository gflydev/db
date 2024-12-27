package db

import (
	"database/sql"
	"github.com/gflydev/core/utils"
	"github.com/jiveio/fluentsql"
	"github.com/jmoiron/sqlx"
	"log"
)

// ========================================================================================
//                                         DB Model
// ========================================================================================

// Raw struct
type Raw struct {
	sqlStr string
	args   []any
}

type DBModel struct {
	tx *sqlx.Tx

	model any // Model struct
	raw   Raw // Raw struct

	selectStatement      fluentsql.Select // Select columns
	omitsSelectStatement fluentsql.Select // Omit columns
	whereStatement       fluentsql.Where  // Where conditions
	joinStatement        fluentsql.Join
	groupByStatement     fluentsql.GroupBy
	havingStatement      fluentsql.Having // A version of Where
	orderByStatement     fluentsql.OrderBy
	limitStatement       fluentsql.Limit
	fetchStatement       fluentsql.Fetch // A version of Limit
}

func Instance() *DBModel {
	return &DBModel{
		tx:    nil,
		model: nil,
	}
}

// Reset DB model's builders after everytime perform the DB query
func (db *DBModel) reset() *DBModel {
	db.model = nil
	db.raw.sqlStr = ""
	db.selectStatement.Columns = []any{}
	db.omitsSelectStatement.Columns = []any{}
	db.whereStatement.Conditions = []fluentsql.Condition{}
	db.joinStatement.Items = []fluentsql.JoinItem{}
	db.groupByStatement.Items = []string{}
	db.havingStatement.Conditions = []fluentsql.Condition{}
	db.orderByStatement.Items = []fluentsql.SortItem{}
	db.limitStatement.Limit = 0
	db.fetchStatement.Fetch = 0

	return db
}

// ========================================================================================
//                                 FluentSQL + SQLX integration
// ========================================================================================

// get perform getting single data row by QueryBuilder
func (db *DBModel) get(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.getRaw(sqlStr, args, model)
}

// get perform getting single data row by QueryBuilder
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

// query performs query list data row by QueryBuilder
func (db *DBModel) query(q *fluentsql.QueryBuilder, model any) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.queryRaw(sqlStr, args, model)
}

// queryRaw performs query list data row by sqlStr and arguments
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

// add performs adding new data by InsertBuilder
func (db *DBModel) add(q *fluentsql.InsertBuilder, primaryColumn string) (id any, err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.addRaw(sqlStr, args, primaryColumn)
}

// addRaw performs adding new data by sqlStr and arguments
func (db *DBModel) addRaw(sqlStr string, args []any, primaryColumn string) (id any, err error) {
	if utils.Getenv("DB_DEBUG", false) {
		log.Printf("SQL> %s - args %v", sqlStr, args)
	}

	// Data persistence
	if fluentsql.DBType() == fluentsql.PostgreSQL {
		if primaryColumn != "" {
			sqlStr += " RETURNING " + primaryColumn

			if utils.Getenv("DB_DEBUG", false) {
				log.Printf("Chagned SQL> %s", sqlStr)
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

// update performs updating data by UpdateBuilder
func (db *DBModel) update(q *fluentsql.UpdateBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// delete performs deleting data by DeleteBuilder
func (db *DBModel) delete(q *fluentsql.DeleteBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// execRaw performs updating and deleting data by DeleteBuilder
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

// Count get total rows
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

// ========================================================================================
//                                 DB Model operators
// ========================================================================================

// Raw build query from raw SQL
func (db *DBModel) Raw(sqlStr string, args ...any) *DBModel {
	db.raw.sqlStr = sqlStr
	db.raw.args = args

	return db
}

// Select List of columns
func (db *DBModel) Select(columns ...any) *DBModel {
	db.selectStatement.Columns = columns

	return db
}

// Omit exclude some columns
func (db *DBModel) Omit(columns ...any) *DBModel {
	db.omitsSelectStatement.Columns = columns

	return db
}

// Model set specific model for builder
func (db *DBModel) Model(model any) *DBModel {
	db.model = model

	return db
}

// Where add where condition
func (db *DBModel) Where(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// WhereOr add where condition
func (db *DBModel) WhereOr(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.whereStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.Or,
	})

	return db
}

// WhereGroup combine multi where conditions into a group.
func (db *DBModel) WhereGroup(groupCondition fluentsql.FnWhereBuilder) *DBModel {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	cond := fluentsql.Condition{
		Group: whereBuilder.Conditions(),
	}

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, cond)

	return db
}

// When checking TRUE to build Where condition.
func (db *DBModel) When(condition bool, groupCondition fluentsql.FnWhereBuilder) *DBModel {
	if !condition {
		return db
	}

	// Create new WhereBuilder
	whereBuilder := groupCondition(*fluentsql.WhereInstance())

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, whereBuilder.Conditions()...)

	return db
}

// Join builder
func (db *DBModel) Join(join fluentsql.JoinType, table string, condition fluentsql.Condition) *DBModel {
	db.joinStatement.Append(fluentsql.JoinItem{
		Join:      join,
		Table:     table,
		Condition: condition,
	})

	return db
}

// Having builder
func (db *DBModel) Having(field any, opt fluentsql.WhereOpt, value any) *DBModel {
	db.havingStatement.Append(fluentsql.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: fluentsql.And,
	})

	return db
}

// GroupBy fields in a query
func (db *DBModel) GroupBy(fields ...string) *DBModel {
	db.groupByStatement.Append(fields...)

	return db
}

// OrderBy builder
func (db *DBModel) OrderBy(field string, dir fluentsql.OrderByDir) *DBModel {
	db.orderByStatement.Append(field, dir)

	return db
}

// Limit builder
func (db *DBModel) Limit(limit, offset int) *DBModel {
	db.limitStatement.Limit = limit
	db.limitStatement.Offset = offset

	return db
}

// RemoveLimit builder
func (db *DBModel) RemoveLimit() fluentsql.Limit {
	var _limitStatement fluentsql.Limit

	_limitStatement.Limit = db.limitStatement.Limit
	_limitStatement.Offset = db.limitStatement.Offset

	db.limitStatement.Limit = 0
	db.limitStatement.Offset = 0

	return _limitStatement
}

// Fetch builder
func (db *DBModel) Fetch(offset, fetch int) *DBModel {
	db.fetchStatement.Offset = offset
	db.fetchStatement.Fetch = fetch

	return db
}

// RemoveFetch builder
func (db *DBModel) RemoveFetch() fluentsql.Fetch {
	var _fetchStatement fluentsql.Fetch

	_fetchStatement.Offset = db.fetchStatement.Offset
	_fetchStatement.Fetch = db.fetchStatement.Fetch

	db.fetchStatement.Offset = 0
	db.fetchStatement.Fetch = 0

	return _fetchStatement
}

// whereFromModel Build and append WHERE clause from specific model's data off table.
func (tbl *Table) whereFromModel(queryBuilder *fluentsql.QueryBuilder) {
	if tbl.HasData {
		for _, column := range tbl.Columns {
			// Prevent some meta, relational, and default (Zero) value of column
			if column.isNotData() || column.IsZero {
				continue
			}

			// Append query conditions
			queryBuilder.Where(column.Name, fluentsql.Eq, tbl.Values[column.Name])
		}
	}
}

// ========================================================================================
//                                     DB Transaction
// ========================================================================================

// Begin new transaction
func (db *DBModel) Begin() *DBModel {
	db.tx = dbInstanceTx()

	return db
}

// Rollback transaction
func (db *DBModel) Rollback() error {
	if db.tx != nil {
		return db.tx.Rollback()
	}

	return nil
}

// Commit transaction
func (db *DBModel) Commit() error {
	if db.tx != nil {
		return db.tx.Commit()
	}

	return nil
}
