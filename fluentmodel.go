package db

import (
	"database/sql"
	"github.com/gflydev/core/utils"
	qb "github.com/jivegroup/fluentsql"
	"github.com/jmoiron/sqlx"
	"log"
)

// ====================================================================
//                              DB Model
// ====================================================================

// Raw represents a raw SQL query with its associated parameters for direct database execution.
// This struct encapsulates both the SQL query string and its corresponding arguments,
// providing a way to execute custom SQL queries that bypass the ORM's query builder.
// It's particularly useful for complex queries, database-specific operations, or
// performance-critical scenarios where hand-crafted SQL is preferred.
//
// Fields:
//   - sqlStr (string): The raw SQL query string with parameter placeholders.
//     Supports standard SQL parameter placeholders (? for most databases, $1, $2, etc. for PostgreSQL)
//   - args ([]any): Slice of arguments that correspond to the parameter placeholders in sqlStr.
//     Arguments are applied in order and should match the placeholder count and types
//
// Usage:
//
//	// The Raw struct is typically used internally by the Raw() method
//	raw := Raw{
//	    sqlStr: "SELECT * FROM users WHERE age > ? AND status = ?",
//	    args:   []any{18, "active"},
//	}
//
//	// More commonly used through the fluent interface:
//	db.Raw("SELECT * FROM users WHERE age > ? AND status = ?", 18, "active")
//
// Examples of supported SQL operations:
//   - Complex SELECT queries with subqueries and CTEs
//   - Database-specific functions and operators
//   - Stored procedure calls
//   - Data manipulation with custom logic
//   - Performance-optimized queries
//   - Database administration commands
type Raw struct {
	sqlStr string // The SQL query string with parameter placeholders for safe execution
	args   []any  // The ordered arguments corresponding to query placeholders
}

// DBModel represents the core fluent interface for database operations and query building.
// This struct serves as the primary entry point for all database interactions, providing
// a chainable API for constructing complex SQL queries, managing transactions, and
// executing database operations. It encapsulates various SQL clause builders and
// maintains state for both ORM-style operations and raw SQL execution.
//
// Fields:
//
//   - tx (*sqlx.Tx): Optional database transaction for atomic operations.
//     When set, all database operations will be executed within this transaction context.
//     Nil indicates operations should use the global database connection.
//
//   - model (any): The target model struct that defines the database table structure.
//     Used for ORM operations to determine table name, column mappings, and data types.
//     Should be a struct or pointer to struct with appropriate database tags.
//
//   - raw (Raw): Container for raw SQL queries and their parameters.
//     When populated, takes precedence over query builder operations.
//     Allows execution of custom SQL that bypasses the ORM query construction.
//
//   - selectStatement (qb.Select): Builder for SELECT clause construction.
//     Manages column selection, including specific columns, expressions, and aliases.
//     Used in conjunction with query operations to control result set structure.
//
//   - omitsSelectStatement (qb.Select): Builder for column omission in SELECT operations.
//     Specifies columns to exclude from SELECT statements, useful for excluding
//     sensitive or unnecessary fields from query results.
//
//   - whereStatement (qb.Where): Builder for WHERE clause conditions.
//     Manages filtering conditions including AND, OR, and grouped conditions.
//     Supports various comparison operators and complex conditional logic.
//
//   - joinStatement (qb.Join): Builder for JOIN operations between tables.
//     Handles INNER, LEFT, RIGHT, and FULL OUTER joins with custom conditions.
//     Enables relational queries across multiple tables.
//
//   - groupByStatement (qb.GroupBy): Builder for GROUP BY clause construction.
//     Manages column grouping for aggregate operations and result set organization.
//     Works in conjunction with aggregate functions and HAVING clauses.
//
//   - havingStatement (qb.Having): Builder for HAVING clause conditions.
//     Provides filtering capabilities for grouped results, similar to WHERE
//     but operates on aggregated data after GROUP BY operations.
//
//   - orderByStatement (qb.OrderBy): Builder for ORDER BY clause construction.
//     Manages result set sorting with support for multiple columns and directions.
//     Enables ascending and descending sort orders with custom priority.
//
//   - limitStatement (qb.Limit): Builder for LIMIT and OFFSET operations.
//     Controls result set size and pagination for efficient data retrieval.
//     Supports both limit-only and limit-with-offset configurations.
//
//   - fetchStatement (qb.Fetch): Builder for FETCH clause operations.
//     Alternative to LIMIT for databases that support FETCH FIRST/NEXT syntax.
//     Provides SQL standard-compliant result set limiting.
//
// Usage Patterns:
//
//	// Basic query building
//	db := Instance().Model(&User{}).Where("age", qb.Gt, 18).OrderBy("name", qb.Asc)
//
//	// Transaction-based operations
//	db := Instance().Begin()
//	defer db.Rollback() // Rollback if not committed
//
//	// Raw SQL execution
//	db := Instance().Raw("SELECT * FROM users WHERE custom_condition")
//
//	// Complex query with multiple clauses
//	db := Instance().Model(&User{}).
//	    Select("name", "email").
//	    Where("status", qb.Eq, "active").
//	    Join(qb.InnerJoin, "profiles", qb.Condition{Field: "users.id", Opt: qb.Eq, Value: "profiles.user_id"}).
//	    GroupBy("department").
//	    Having("COUNT(*)", qb.Gt, 5).
//	    OrderBy("created_at", qb.Desc).
//	    Limit(10, 0)
//
// Note:
//   - All builder fields are reset after each operation to prevent state leakage
//   - Transaction field persists across operations until explicitly committed or rolled back
//   - Raw SQL takes precedence over query builder operations when both are present
//   - The struct is designed for method chaining to create fluent, readable database code
type DBModel struct {
	tx *sqlx.Tx // Database transaction context for atomic operations

	model any // Target model struct defining table structure and column mappings
	raw   Raw // Raw SQL query container with parameters for custom query execution

	selectStatement      qb.Select  // SELECT clause builder for column specification and result shaping
	omitsSelectStatement qb.Select  // Column omission builder for excluding specific fields from results
	whereStatement       qb.Where   // WHERE clause builder for filtering conditions and logical operations
	joinStatement        qb.Join    // JOIN clause builder for multi-table relational queries
	groupByStatement     qb.GroupBy // GROUP BY clause builder for result aggregation and organization
	havingStatement      qb.Having  // HAVING clause builder for post-aggregation filtering
	orderByStatement     qb.OrderBy // ORDER BY clause builder for result sorting and ordering
	limitStatement       qb.Limit   // LIMIT clause builder for result set size control and pagination
	fetchStatement       qb.Fetch   // FETCH clause builder for SQL standard-compliant result limiting
}

// Instance creates and returns a new DBModel instance for database operations.
// This function serves as the primary entry point for creating database query builders
// and performing ORM operations. Each instance is independent and maintains its own
// state for query building, allowing for concurrent usage and isolated operations.
// The returned instance provides a fluent interface for chaining database operations.
//
// Returns:
//   - *DBModel: A new database model instance with:
//   - Clean state with no active transaction
//   - Empty query builders ready for configuration
//   - No associated model (must be set via Model() method)
//   - All SQL clause builders initialized to empty state
//
// Examples:
//
//	// Basic instance creation for simple queries
//	db := Instance()
//	users := []User{}
//	err := db.Model(&User{}).Find(&users)
//
//	// Chained operations with fluent interface
//	db := Instance().Model(&User{}).Where("age", qb.Gt, 18).OrderBy("name", qb.Asc)
//
//	// Multiple independent instances
//	userDB := Instance().Model(&User{})
//	productDB := Instance().Model(&Product{})
//	// Each instance maintains separate state
//
//	// Transaction-based operations
//	db := Instance().Begin()
//	defer func() {
//	    if r := recover(); r != nil {
//	        db.Rollback()
//	    }
//	}()
//
//	// Perform operations within transaction
//	err := db.Model(&User{}).Create(user)
//	if err != nil {
//	    db.Rollback()
//	    return err
//	}
//
//	return db.Commit()
//
//	// Raw SQL operations
//	db := Instance().Raw("SELECT COUNT(*) FROM users WHERE status = ?", "active")
//	var count int
//	err := db.Get(&count)
//
// Usage Patterns:
//   - Create instance → Set model → Build query → Execute operation
//   - Instance().Model(&Model{}).QueryMethod().ExecutionMethod()
//   - Each instance is independent and thread-safe for its own operations
//   - Instances can be reused but are automatically reset after operations
//
// Note:
//   - Each call creates a completely new instance with clean state
//   - No database connection is established until an operation is executed
//   - The instance must be configured with a model before most operations
//   - Transaction state persists across operations until committed or rolled back
//   - All query builder state is reset after each terminal operation
func Instance() *DBModel {
	return &DBModel{
		tx:    nil, // No active transaction initially
		model: nil, // No model associated initially, must be set via Model()
	}
}

// reset clears the state of the DBModel and resets builders.
//
// Returns:
//
//	*DBModel - The reset DBModel instance.
func (db *DBModel) reset() *DBModel {
	db.model = nil                                   // Clear the model.
	db.raw.sqlStr = ""                               // Reset raw SQL string.
	db.selectStatement.Columns = []any{}             // Clear SELECT columns.
	db.omitsSelectStatement.Columns = []any{}        // Clear omitted SELECT columns.
	db.whereStatement.Conditions = []qb.Condition{}  // Clear WHERE conditions.
	db.joinStatement.Items = []qb.JoinItem{}         // Clear JOIN items.
	db.groupByStatement.Items = []string{}           // Clear GROUP BY items.
	db.havingStatement.Conditions = []qb.Condition{} // Clear HAVING conditions.
	db.orderByStatement.Items = []qb.SortItem{}      // Clear ORDER BY items.
	db.limitStatement.Limit = 0                      // Reset limit.
	db.fetchStatement.Fetch = 0                      // Reset fetch.

	return db
}

// ====================================================================
//                      FluentSQL + SQLX integration
// ====================================================================

// get performs fetching a single data row using QueryBuilder.
//
// Parameters:
//   - q (*qb.QueryBuilder): The query builder comprising the SQL query and arguments.
//   - model (any): The model to map the resulting row.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) get(q *qb.QueryBuilder, model any) (err error) {
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
//   - q (*qb.QueryBuilder): The query builder with the SQL and arguments.
//   - model (any): The model to map the resulting rows.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) query(q *qb.QueryBuilder, model any) (err error) {
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
//   - q (*qb.InsertBuilder): The insert query builder with the SQL and arguments.
//   - primaryColumn (string): The primary column to return, used for PostgreSQL.
//
// Returns:
//   - id (any): The ID of the newly inserted row.
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) add(q *qb.InsertBuilder, primaryColumn string) (id any, err error) {
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
	if qb.IsDialect(qb.PostgreSQL) {
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
	} else if qb.IsDialect(qb.MySQL) {
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
//   - q (*qb.UpdateBuilder): The update query builder with the SQL and arguments.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) update(q *qb.UpdateBuilder) (err error) {
	var sqlStr string
	var args []any

	sqlStr, args, _ = q.Sql()

	return db.execRaw(sqlStr, args)
}

// delete performs deleting data using DeleteBuilder.
//
// Parameters:
//   - q (*qb.DeleteBuilder): The delete query builder with the SQL and arguments.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) delete(q *qb.DeleteBuilder) (err error) {
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
//   - q (*qb.QueryBuilder): The query builder with the SQL and arguments.
//   - total (*int): Pointer to an integer to store the total count.
//
// Returns:
//   - err (error): Error encountered during execution, if any.
func (db *DBModel) count(q *qb.QueryBuilder, total *int) error {
	var fetch qb.Fetch
	var limit qb.Limit

	// Build SQL without pagination
	fetch = q.RemoveFetch()
	limit = q.RemoveLimit()

	// Create COUNT query
	sqlBuilderCount := qb.QueryInstance().
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

// Raw configures the DBModel to execute a custom SQL query with parameters.
// This method allows bypassing the ORM's query builder to execute hand-crafted SQL
// for complex operations, database-specific features, or performance optimization.
// When Raw is used, it takes precedence over all other query building methods,
// and the provided SQL will be executed directly against the database.
//
// Parameters:
//   - sqlStr (string): The raw SQL query string with parameter placeholders.
//     Supports database-specific placeholder syntax:
//   - MySQL/SQLite: Use ? for parameters (e.g., "SELECT * FROM users WHERE id = ?")
//   - PostgreSQL: Use $1, $2, etc. for parameters (e.g., "SELECT * FROM users WHERE id = $1")
//   - args (...any): Variadic arguments that correspond to the parameter placeholders.
//     Arguments are applied in order and should match the placeholder count and expected types.
//     Supports all Go types that can be converted to SQL types.
//
// Returns:
//   - *DBModel: The same DBModel instance for method chaining, configured for raw SQL execution.
//
// Examples:
//
//	// Simple SELECT with parameters
//	var users []User
//	err := Instance().Raw("SELECT * FROM users WHERE age > ? AND status = ?", 18, "active").Find(&users)
//
//	// Complex query with joins and subqueries
//	var results []UserStats
//	query := `
//	    SELECT u.name, COUNT(o.id) as order_count, AVG(o.total) as avg_order
//	    FROM users u
//	    LEFT JOIN orders o ON u.id = o.user_id
//	    WHERE u.created_at > ?
//	    GROUP BY u.id, u.name
//	    HAVING COUNT(o.id) > ?
//	    ORDER BY avg_order DESC
//	`
//	err := Instance().Raw(query, time.Now().AddDate(-1, 0, 0), 5).Find(&results)
//
//	// Database-specific functions
//	var count int
//	err := Instance().Raw("SELECT COUNT(*) FROM users WHERE MATCH(name) AGAINST(?)", "john").Get(&count)
//
//	// INSERT with RETURNING (PostgreSQL)
//	var newID int
//	err := Instance().Raw("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", "John", "john@example.com").Get(&newID)
//
//	// UPDATE with complex conditions
//	err := Instance().Raw("UPDATE products SET price = price * 0.9 WHERE category_id IN (SELECT id FROM categories WHERE name = ?)", "electronics").Exec()
//
//	// DELETE with joins
//	err := Instance().Raw("DELETE u FROM users u JOIN user_sessions s ON u.id = s.user_id WHERE s.last_activity < ?", oldDate).Exec()
//
//	// Stored procedure call
//	err := Instance().Raw("CALL update_user_statistics(?)", userID).Exec()
//
//	// Common Table Expression (CTE)
//	query := `
//	    WITH RECURSIVE category_tree AS (
//	        SELECT id, name, parent_id, 0 as level FROM categories WHERE parent_id IS NULL
//	        UNION ALL
//	        SELECT c.id, c.name, c.parent_id, ct.level + 1
//	        FROM categories c
//	        JOIN category_tree ct ON c.parent_id = ct.id
//	    )
//	    SELECT * FROM category_tree WHERE level <= ?
//	`
//	var categories []Category
//	err := Instance().Raw(query, 3).Find(&categories)
//
// Use Cases:
//   - Complex analytical queries with window functions
//   - Database-specific optimizations and hints
//   - Stored procedure and function calls
//   - Bulk operations with custom logic
//   - Migration and schema manipulation queries
//   - Performance-critical queries requiring manual optimization
//
// Note:
//   - Raw SQL takes precedence over all query builder methods
//   - Parameter placeholders prevent SQL injection when used properly
//   - The query is executed in the current transaction context if one exists
//   - Database-specific syntax should match the configured database driver
//   - Use appropriate execution methods: Get() for single row, Find() for multiple rows, Exec() for non-query operations
//   - SQL debugging can be enabled via DB_DEBUG environment variable
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

// Model associates a struct model with the DBModel instance for ORM operations.
// This method is fundamental to the ORM functionality as it defines the target table
// structure, column mappings, and data types for database operations. The model
// struct should have appropriate database tags to define table and column mappings.
// This method must be called before most ORM operations to establish the context.
//
// Parameters:
//   - model (any): The model struct that defines the database table structure. Supported types:
//   - Struct: Direct struct value (e.g., User{})
//   - *Struct: Pointer to struct (e.g., &User{}) - recommended for most operations
//   - The struct should have database tags for proper column mapping:
//   - `db:"column_name"` for column name mapping
//   - `db:"column_name,primary"` for primary key designation
//   - `db:"-"` to exclude fields from database operations
//
// Returns:
//   - *DBModel: The same DBModel instance for method chaining, now configured with the model.
//
// Examples:
//
//	// Basic model association
//	type User struct {
//	    ID    int    `db:"id,primary"`
//	    Name  string `db:"name"`
//	    Email string `db:"email"`
//	    Age   int    `db:"age"`
//	}
//
//	db := Instance().Model(&User{})
//
//	// Query operations after model association
//	var users []User
//	err := db.Find(&users)
//
//	// Create operations
//	user := User{Name: "John", Email: "john@example.com", Age: 30}
//	err := db.Create(&user) // user.ID will be populated
//
//	// Update operations
//	user.Age = 31
//	err := db.Where("id", Eq, user.ID).Update(&user)
//
//	// Delete operations
//	err := db.Where("id", Eq, user.ID).Delete(&user)
//
//	// Complex model with relationships
//	type Order struct {
//	    ID       int       `db:"id,primary"`
//	    UserID   int       `db:"user_id"`
//	    Total    float64   `db:"total"`
//	    Status   string    `db:"status"`
//	    Created  time.Time `db:"created_at"`
//	    Updated  time.Time `db:"updated_at"`
//	    Internal string    `db:"-"` // Excluded from database operations
//	}
//
//	db := Instance().Model(&Order{})
//
//	// Model with custom table name (if different from struct name)
//	type UserProfile struct {
//	    ID     int    `db:"id,primary"`
//	    UserID int    `db:"user_id"`
//	    Bio    string `db:"biography"`
//	} // Maps to "user_profiles" table by convention
//
//	// Multiple models for different operations
//	userDB := Instance().Model(&User{})
//	orderDB := Instance().Model(&Order{})
//	// Each instance maintains separate model context
//
//	// Model reuse with different instances
//	db1 := Instance().Model(&User{}).Where("status", Eq, "active")
//	db2 := Instance().Model(&User{}).Where("age", qb.Gt, 18)
//	// Independent query contexts with same model
//
// Model Requirements:
//   - Struct fields should be exported (capitalized)
//   - Use appropriate database tags for column mapping
//   - Primary key fields should be tagged with "primary"
//   - Time fields should use time.Time type for proper handling
//   - Nullable fields can use sql.NullString, sql.NullInt64, etc.
//
// Common Patterns:
//   - Always use pointer to struct: Model(&User{}) instead of Model(User{})
//   - Set model before any ORM operations
//   - One model per DBModel instance for clarity
//   - Use consistent naming conventions between struct and database
//
// Note:
//   - The model defines the target table name (derived from struct name)
//   - Column mappings are established through struct tags
//   - Primary key detection is automatic based on "primary" tag
//   - The model persists until the DBModel instance is reset
//   - Required for Create, Update, Delete, Find, and other ORM operations
//   - Raw SQL operations can work without a model but lose ORM benefits
func (db *DBModel) Model(model any) *DBModel {
	db.model = model

	return db
}

// Where adds a WHERE condition to the query.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (qb.WhereOpt): The operator to use (e.g., equals, greater than).
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Where(field any, opt WhereOpt, value any) *DBModel {
	db.whereStatement.Append(qb.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
	})

	return db
}

// WhereOr adds an OR condition to the WHERE clause.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (qb.WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) WhereOr(field any, opt WhereOpt, value any) *DBModel {
	db.whereStatement.Append(qb.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: qb.Or,
	})

	return db
}

// WhereGroup combines multiple WHERE conditions into a group.
//
// Parameters:
//   - groupCondition (qb.FnWhereBuilder): The function to build grouped conditions.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) WhereGroup(groupCondition qb.FnWhereBuilder) *DBModel {
	// Create new WhereBuilder
	whereBuilder := groupCondition(*qb.WhereInstance())

	cond := qb.Condition{
		Group: whereBuilder.Conditions(),
	}

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, cond)

	return db
}

// When conditionally applies a WHERE condition if the provided condition is TRUE.
//
// Parameters:
//   - condition (bool): Determines whether the condition should be applied.
//   - groupCondition (qb.FnWhereBuilder): The function to build the condition.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) When(condition bool, groupCondition qb.FnWhereBuilder) *DBModel {
	if !condition {
		return db
	}

	// Create new WhereBuilder
	whereBuilder := groupCondition(*qb.WhereInstance())

	db.whereStatement.Conditions = append(db.whereStatement.Conditions, whereBuilder.Conditions()...)

	return db
}

// Join adds a JOIN clause to the query.
//
// Parameters:
//   - join (qb.JoinType): The type of join (e.g., INNER, LEFT).
//   - table (string): The name of the table to join.
//   - condition (qb.Condition): The condition for the join.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Join(join qb.JoinType, table string, condition Condition) *DBModel {
	db.joinStatement.Append(qb.JoinItem{
		Join:      join,
		Table:     table,
		Condition: condition.ToQBCondition(),
	})

	return db
}

// Having adds a HAVING condition to the query.
//
// Parameters:
//   - field (any): The field or column to filter.
//   - opt (qb.WhereOpt): The operator to use.
//   - value (any): The value to compare against.
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) Having(field any, opt WhereOpt, value any) *DBModel {
	db.havingStatement.Append(qb.Condition{
		Field: field,
		Opt:   opt,
		Value: value,
		AndOr: And,
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
//   - dir (qb.OrderByDir): The sorting direction (e.g., ASC or DESC).
//
// Returns:
//   - *DBModel: A reference to the DBModel instance for chaining.
func (db *DBModel) OrderBy(field string, dir qb.OrderByDir) *DBModel {
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
//   - qb.Limit: The removed limit settings.
func (db *DBModel) RemoveLimit() qb.Limit {
	var _limitStatement qb.Limit

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
//   - qb.Fetch: The removed fetch settings.
func (db *DBModel) RemoveFetch() qb.Fetch {
	var _fetchStatement qb.Fetch

	_fetchStatement.Offset = db.fetchStatement.Offset
	_fetchStatement.Fetch = db.fetchStatement.Fetch

	db.fetchStatement.Offset = 0
	db.fetchStatement.Fetch = 0

	return _fetchStatement
}

// whereFromModel builds and appends a WHERE clause from the specific model's data.
//
// Parameters:
//   - queryBuilder (*qb.QueryBuilder): The query builder to modify.
func (tbl *Table) whereFromModel(queryBuilder *qb.QueryBuilder) {
	if tbl.HasData {
		for _, column := range tbl.Columns {
			// Prevent processing meta, relational, and default (zero) column values
			if column.isNotData() || column.IsZero {
				continue
			}

			// Append query conditions
			queryBuilder.Where(column.Name, Eq, tbl.Values[column.Name])
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
