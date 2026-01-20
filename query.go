// Package sqlquery provides a fluent API for building SQL queries with support for MySQL, PostgreSQL, and SQLite.
//
// The package wraps github.com/huandu/go-sqlbuilder and provides simplified query building functions
// with support for advanced filtering, pagination, and row locking.
//
// # Basic Usage
//
// Build a SELECT query with filters:
//
//	options := NewFindAllOptions(MySQLFlavor).
//		WithFilter("status", "active").
//		WithFilter("age.gte", 18).
//		WithLimit(10).
//		WithOffset(0).
//		WithOrderBy("created_at DESC")
//	sql, args := FindAllQuery("users", options)
//
// # Filter Syntax
//
// The Filters map supports special operators via dot notation:
//
//	"field"          - Equality (field = value)
//	"field.in"       - IN clause (value must be comma-separated string)
//	"field.notin"    - NOT IN clause (value must be comma-separated string)
//	"field.not"      - Not equal (field != value)
//	"field.gt"       - Greater than (field > value)
//	"field.gte"      - Greater or equal (field >= value)
//	"field.lt"       - Less than (field < value)
//	"field.lte"      - Less or equal (field <= value)
//	"field.like"     - LIKE pattern matching
//	"field.null"     - IS NULL / IS NOT NULL (value must be bool)
//
// # Supported Database Flavors
//
// Use one of the predefined flavors when creating options:
//
//	MySQLFlavor      - MySQL and MariaDB
//	PostgreSQLFlavor - PostgreSQL
//	SQLiteFlavor     - SQLite
package sqlquery

import (
	"sort"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

// parseIn converts a comma-separated string into a slice of interface{} values
// for use in SQL IN clauses. For example, "1,2,3" becomes []interface{}{"1", "2", "3"}.
func parseIn(value string) []interface{} {
	values := strings.Split(value, ",")
	result := make([]interface{}, len(values))
	for i := range values {
		result[i] = values[i]
	}
	return result
}

// parseSelectFilter applies WHERE conditions to a SELECT query builder based on the filter key and value.
// It supports special filter operators via dot notation (e.g., "field.in", "field.gt", "field.like").
// See package documentation for the full list of supported operators.
func parseSelectFilter(sb *sqlbuilder.SelectBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.NotIn(parsedKey, values...))
			}
		case "not":
			sb.Where(sb.NotEqual(parsedKey, value))
		case "gt":
			sb.Where(sb.GreaterThan(parsedKey, value))
		case "gte":
			sb.Where(sb.GreaterEqualThan(parsedKey, value))
		case "lt":
			sb.Where(sb.LessThan(parsedKey, value))
		case "lte":
			sb.Where(sb.LessEqualThan(parsedKey, value))
		case "like":
			sb.Where(sb.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					sb.Where(sb.IsNull(key))
				} else {
					sb.Where(sb.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			sb.Where(sb.IsNull(key))
		default:
			sb.Where(sb.Equal(key, value))
		}
	}
}

// parseUpdateFilter applies WHERE conditions to an UPDATE query builder based on the filter key and value.
// It supports special filter operators via dot notation (e.g., "field.in", "field.gt", "field.like").
// See package documentation for the full list of supported operators.
func parseUpdateFilter(ub *sqlbuilder.UpdateBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				ub.Where(ub.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				ub.Where(ub.NotIn(parsedKey, values...))
			}
		case "not":
			ub.Where(ub.NotEqual(parsedKey, value))
		case "gt":
			ub.Where(ub.GreaterThan(parsedKey, value))
		case "gte":
			ub.Where(ub.GreaterEqualThan(parsedKey, value))
		case "lt":
			ub.Where(ub.LessThan(parsedKey, value))
		case "lte":
			ub.Where(ub.LessEqualThan(parsedKey, value))
		case "like":
			ub.Where(ub.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					ub.Where(ub.IsNull(key))
				} else {
					ub.Where(ub.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			ub.Where(ub.IsNull(key))
		default:
			ub.Where(ub.Equal(key, value))
		}
	}
}

// parseDeleteFilter applies WHERE conditions to a DELETE query builder based on the filter key and value.
// It supports special filter operators via dot notation (e.g., "field.in", "field.gt", "field.like").
// See package documentation for the full list of supported operators.
func parseDeleteFilter(db *sqlbuilder.DeleteBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				db.Where(db.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				db.Where(db.NotIn(parsedKey, values...))
			}
		case "not":
			db.Where(db.NotEqual(parsedKey, value))
		case "gt":
			db.Where(db.GreaterThan(parsedKey, value))
		case "gte":
			db.Where(db.GreaterEqualThan(parsedKey, value))
		case "lt":
			db.Where(db.LessThan(parsedKey, value))
		case "lte":
			db.Where(db.LessEqualThan(parsedKey, value))
		case "like":
			db.Where(db.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					db.Where(db.IsNull(key))
				} else {
					db.Where(db.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			db.Where(db.IsNull(key))
		default:
			db.Where(db.Equal(key, value))
		}
	}
}

// FindQuery builds a SELECT query and returns the compiled SQL string and arguments.
//
// The function supports filtering with special operators, field selection, and row locking (FOR UPDATE).
//
// Parameters:
//   - tableName: The name of the table to query
//   - options: Configuration including flavor, fields, filters, and FOR UPDATE settings
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	options := NewFindOptions(MySQLFlavor).
//		WithFields([]string{"id", "name", "email"}).
//		WithFilter("status", "active").
//		WithFilter("age.gte", 18)
//	sql, args := FindQuery("users", options)
//	// SELECT id, name, email FROM users WHERE status = ? AND age >= ?
func FindQuery(tableName string, options *FindOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName)
	for key, value := range options.Filters {
		parseSelectFilter(sb, key, value)
	}
	if options.ForUpdate {
		sb.ForUpdate()
		if options.ForUpdateMode != "" {
			sb.SQL(options.ForUpdateMode)
		}
	}
	return sb.Build()
}

// FindAllQuery builds a SELECT query with pagination and returns the compiled SQL string and arguments.
//
// This function extends FindQuery with support for LIMIT, OFFSET, and ORDER BY clauses.
//
// Parameters:
//   - tableName: The name of the table to query
//   - options: Configuration including flavor, fields, filters, limit, offset, ordering, and FOR UPDATE settings
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	options := NewFindAllOptions(MySQLFlavor).
//		WithFilter("status", "active").
//		WithLimit(10).
//		WithOffset(20).
//		WithOrderBy("created_at DESC")
//	sql, args := FindAllQuery("users", options)
//	// SELECT * FROM users WHERE status = ? ORDER BY created_at DESC LIMIT 10 OFFSET 20
func FindAllQuery(tableName string, options *FindAllOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName).Limit(options.Limit).Offset(options.Offset)
	for key, value := range options.Filters {
		parseSelectFilter(sb, key, value)
	}
	if options.OrderBy != "" {
		sb.OrderBy(options.OrderBy)
	}
	if options.ForUpdate {
		sb.ForUpdate()
		if options.ForUpdateMode != "" {
			sb.SQL(options.ForUpdateMode)
		}
	}
	return sb.Build()
}

// InsertQuery builds an INSERT query from a struct and returns the compiled SQL string and arguments.
//
// The function uses struct tags to map struct fields to database columns. Fields are extracted
// based on the specified tag name (e.g., "db", "sql", or custom tags).
//
// Parameters:
//   - flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - tag: The struct tag name to use for field mapping (e.g., "db", "sql")
//   - tableName: The name of the table to insert into
//   - structValue: The struct instance containing values to insert
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	type User struct {
//		Name  string `db:"name"`
//		Email string `db:"email"`
//		Age   int    `db:"age"`
//	}
//	user := User{Name: "John", Email: "john@example.com", Age: 30}
//	sql, args := InsertQuery(MySQLFlavor, "db", "users", user)
//	// INSERT INTO users (name, email, age) VALUES (?, ?, ?)
func InsertQuery(flavor Flavor, tag, tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ib := theStruct.WithTag(tag).InsertInto(tableName, structValue)
	return ib.Build()
}

// UpdateQuery builds an UPDATE query from a struct and returns the compiled SQL string and arguments.
//
// The function updates a record identified by its ID. Struct tags determine which fields to update.
//
// Parameters:
//   - flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - tag: The struct tag name to use for field mapping (e.g., "db", "sql")
//   - tableName: The name of the table to update
//   - id: The ID value for the WHERE id = ? condition
//   - structValue: The struct instance containing updated values
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	type User struct {
//		Name  string `db:"name"`
//		Email string `db:"email"`
//		Age   int    `db:"age"`
//	}
//	user := User{Name: "John", Email: "john@example.com", Age: 31}
//	sql, args := UpdateQuery(MySQLFlavor, "db", "users", 123, user)
//	// UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?
func UpdateQuery(flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ub := theStruct.WithTag(tag).Update(tableName, structValue)
	ub.Where(ub.Equal("id", id))
	return ub.Build()
}

// DeleteQuery builds a DELETE query and returns the compiled SQL string and arguments.
//
// The function deletes a record identified by its ID.
//
// Parameters:
//   - flavor: The SQL dialect (MySQLFlavor, PostgreSQLFlavor, or SQLiteFlavor)
//   - tableName: The name of the table to delete from
//   - id: The ID value for the WHERE id = ? condition
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	sql, args := DeleteQuery(MySQLFlavor, "users", 123)
//	// DELETE FROM users WHERE id = ?
func DeleteQuery(flavor Flavor, tableName string, id interface{}) (string, []interface{}) {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.Flavor(flavor))
	db.DeleteFrom(tableName)
	db.Where(db.Equal("id", id))
	return db.Build()
}

// UpdateWithOptionsQuery builds an UPDATE query with custom field assignments and filters.
//
// This function provides more flexibility than UpdateQuery by allowing:
//   - Custom field assignments (not limited to struct fields)
//   - Multiple WHERE conditions with filter operators
//   - Updating multiple records at once
//
// Parameters:
//   - tableName: The name of the table to update
//   - options: Configuration including flavor, field assignments, and filter conditions
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	options := NewUpdateOptions(MySQLFlavor).
//		WithAssignment("status", "inactive").
//		WithAssignment("updated_at", time.Now()).
//		WithFilter("last_login.lt", time.Now().AddDate(0, -6, 0))
//	sql, args := UpdateWithOptionsQuery("users", options)
//	// UPDATE users SET status = ?, updated_at = ? WHERE last_login < ?
func UpdateWithOptionsQuery(tableName string, options *UpdateOptions) (string, []interface{}) {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	ub.Update(tableName)
	var assignments []string
	for key, value := range options.Assignments {
		assignments = append(assignments, ub.Assign(key, value))
	}
	sort.Strings(assignments)
	ub = ub.Set(assignments...)
	for key, value := range options.Filters {
		parseUpdateFilter(ub, key, value)
	}
	return ub.Build()
}

// DeleteWithOptionsQuery builds a DELETE query with custom filter conditions.
//
// This function provides more flexibility than DeleteQuery by allowing:
//   - Multiple WHERE conditions with filter operators
//   - Deleting multiple records matching complex criteria
//
// Parameters:
//   - tableName: The name of the table to delete from
//   - options: Configuration including flavor and filter conditions
//
// Returns the compiled SQL string and a slice of arguments for parameterized queries.
//
// Example:
//
//	options := NewDeleteOptions(MySQLFlavor).
//		WithFilter("status", "inactive").
//		WithFilter("created_at.lt", time.Now().AddDate(-1, 0, 0))
//	sql, args := DeleteWithOptionsQuery("users", options)
//	// DELETE FROM users WHERE status = ? AND created_at < ?
func DeleteWithOptionsQuery(tableName string, options *DeleteOptions) (string, []interface{}) {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	db.DeleteFrom(tableName)
	for key, value := range options.Filters {
		parseDeleteFilter(db, key, value)
	}
	return db.Build()
}
