# sqlquery

[![Build Status](https://github.com/allisson/sqlquery/workflows/Release/badge.svg)](https://github.com/allisson/sqlquery/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/sqlquery)](https://goreportcard.com/report/github.com/allisson/sqlquery)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/sqlquery)

A fluent, type-safe SQL query builder for Go with support for **MySQL**, **PostgreSQL**, and **SQLite**. Built on top of [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder), sqlquery provides a simplified API for constructing SQL queries with advanced filtering, pagination, and row locking capabilities.

## Features

- **Multiple Database Support** - Works with MySQL, PostgreSQL, and SQLite
- **Advanced Filtering** - Rich filter syntax with operators (IN, NOT IN, >, >=, <, <=, LIKE, IS NULL, etc.)
- **Fluent API** - Method chaining for clean, readable query construction
- **Type-Safe** - Strongly typed options prevent common SQL building errors
- **Struct Mapping** - Build INSERT/UPDATE queries directly from Go structs
- **Pagination Support** - Built-in LIMIT and OFFSET with ordering
- **Row Locking** - Support for FOR UPDATE with various locking modes
- **Parameterized Queries** - Returns SQL with placeholders and args to prevent SQL injection

## Installation

```bash
go get github.com/allisson/sqlquery
```

## Quick Start

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/allisson/sqlquery"
    _ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
    db, _ := sql.Open("postgres", "postgresql://user:pass@localhost/dbname?sslmode=disable")
    defer db.Close()

    // Build a simple SELECT query
    options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
        WithFilter("status", "active").
        WithLimit(10).
        WithOrderBy("created_at DESC")
    
    sql, args := sqlquery.FindAllQuery("users", options)
    
    rows, err := db.Query(sql, args...)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // Process results...
}
```

## Table of Contents

- [Basic Usage](#basic-usage)
  - [SELECT Queries](#select-queries)
  - [INSERT Queries](#insert-queries)
  - [UPDATE Queries](#update-queries)
  - [DELETE Queries](#delete-queries)
- [Advanced Filtering](#advanced-filtering)
- [Pagination and Ordering](#pagination-and-ordering)
- [Row Locking](#row-locking)
- [Struct-Based Queries](#struct-based-queries)
- [Database Flavors](#database-flavors)
- [Complete Examples](#complete-examples)

## Basic Usage

### SELECT Queries

#### Simple SELECT with FindQuery

```go
// SELECT specific fields
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFields([]string{"id", "name", "email"}).
    WithFilter("status", "active")

sql, args := sqlquery.FindQuery("users", options)
// SQL: SELECT id, name, email FROM users WHERE status = ?
// Args: ["active"]
```

#### Paginated SELECT with FindAllQuery

```go
// SELECT with pagination
options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("status", "active").
    WithLimit(50).
    WithOffset(100).
    WithOrderBy("created_at DESC")

sql, args := sqlquery.FindAllQuery("users", options)
// SQL: SELECT * FROM users WHERE status = $1 ORDER BY created_at DESC LIMIT 50 OFFSET 100
// Args: ["active"]
```

### INSERT Queries

```go
type User struct {
    ID        int    `db:"id"`
    Name      string `db:"name"`
    Email     string `db:"email"`
    Status    string `db:"status"`
}

user := User{
    ID:     1,
    Name:   "John Doe",
    Email:  "john@example.com",
    Status: "active",
}

sql, args := sqlquery.InsertQuery(
    sqlquery.MySQLFlavor,
    "db",        // struct tag to use
    "users",     // table name
    &user,       // struct with values
)
// SQL: INSERT INTO users (id, name, email, status) VALUES (?, ?, ?, ?)
// Args: [1, "John Doe", "john@example.com", "active"]
```

### UPDATE Queries

#### Update by ID using struct

```go
type User struct {
    ID    int    `db:"id"`
    Name  string `db:"name" fieldtag:"update"`
    Email string `db:"email" fieldtag:"update"`
}

user := User{
    ID:    1,
    Name:  "Jane Doe",
    Email: "jane@example.com",
}

sql, args := sqlquery.UpdateQuery(
    sqlquery.PostgreSQLFlavor,
    "update",    // use "fieldtag" to determine which fields to update
    "users",
    user.ID,
    &user,
)
// SQL: UPDATE users SET name = $1, email = $2 WHERE id = $3
// Args: ["Jane Doe", "jane@example.com", 1]
```

#### Update with custom options

```go
options := sqlquery.NewUpdateOptions(sqlquery.MySQLFlavor).
    WithAssignment("status", "inactive").
    WithAssignment("deactivated_at", time.Now()).
    WithFilter("last_login.lt", time.Now().AddDate(0, -6, 0))

sql, args := sqlquery.UpdateWithOptionsQuery("users", options)
// SQL: UPDATE users SET deactivated_at = ?, status = ? WHERE last_login < ?
// Args: [<timestamp>, "inactive", <6 months ago>]
```

### DELETE Queries

#### Delete by ID

```go
sql, args := sqlquery.DeleteQuery(sqlquery.SQLiteFlavor, "users", 1)
// SQL: DELETE FROM users WHERE id = ?
// Args: [1]
```

#### Delete with filters

```go
options := sqlquery.NewDeleteOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("status", "inactive").
    WithFilter("created_at.lt", time.Now().AddDate(-1, 0, 0))

sql, args := sqlquery.DeleteWithOptionsQuery("users", options)
// SQL: DELETE FROM users WHERE created_at < $1 AND status = $2
// Args: [<1 year ago>, "inactive"]
```

## Advanced Filtering

sqlquery supports a rich filter syntax using dot notation for operators:

| Filter Syntax | SQL Operator | Example |
|--------------|--------------|---------|
| `field` | `=` | `WithFilter("status", "active")` → `status = ?` |
| `field.in` | `IN` | `WithFilter("id.in", "1,2,3")` → `id IN (?, ?, ?)` |
| `field.notin` | `NOT IN` | `WithFilter("id.notin", "1,2,3")` → `id NOT IN (?, ?, ?)` |
| `field.not` | `<>` | `WithFilter("status.not", "deleted")` → `status <> ?` |
| `field.gt` | `>` | `WithFilter("age.gt", 18)` → `age > ?` |
| `field.gte` | `>=` | `WithFilter("age.gte", 18)` → `age >= ?` |
| `field.lt` | `<` | `WithFilter("age.lt", 65)` → `age < ?` |
| `field.lte` | `<=` | `WithFilter("age.lte", 65)` → `age <= ?` |
| `field.like` | `LIKE` | `WithFilter("name.like", "%John%")` → `name LIKE ?` |
| `field.null` | `IS NULL` / `IS NOT NULL` | `WithFilter("deleted_at.null", true)` → `deleted_at IS NULL` |
| `field` (nil value) | `IS NULL` | `WithFilter("deleted_at", nil)` → `deleted_at IS NULL` |

### Filter Examples

```go
// Equality filter
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("status", "active")
// WHERE status = ?

// IN filter (comma-separated string)
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("id.in", "1,2,3,4,5")
// WHERE id IN (?, ?, ?, ?, ?)

// NOT IN filter
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("status.notin", "deleted,archived")
// WHERE status NOT IN (?, ?)

// Comparison filters
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("age.gte", 18).
    WithFilter("age.lt", 65)
// WHERE age >= ? AND age < ?

// LIKE filter
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("email.like", "%@example.com")
// WHERE email LIKE ?

// NULL filters
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("deleted_at.null", true)
// WHERE deleted_at IS NULL

options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("verified_at.null", false)
// WHERE verified_at IS NOT NULL

// Multiple filters combined
options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("status", "active").
    WithFilter("age.gte", 18).
    WithFilter("role.in", "admin,moderator").
    WithFilter("deleted_at.null", true).
    WithLimit(10)
// WHERE status = $1 AND age >= $2 AND role IN ($3, $4) AND deleted_at IS NULL LIMIT 10
```

## Pagination and Ordering

```go
// Basic pagination
options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
    WithLimit(20).
    WithOffset(0)

sql, args := sqlquery.FindAllQuery("products", options)
// SELECT * FROM products LIMIT 20 OFFSET 0

// Pagination with ordering
options := sqlquery.NewFindAllOptions(sqlquery.MySQLFlavor).
    WithLimit(50).
    WithOffset(100).
    WithOrderBy("created_at DESC, id ASC")

sql, args := sqlquery.FindAllQuery("orders", options)
// SELECT * FROM orders ORDER BY created_at DESC, id ASC LIMIT 50 OFFSET 100

// Combining with filters
options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("status", "pending").
    WithFilter("priority.gte", 5).
    WithOrderBy("priority DESC, created_at ASC").
    WithLimit(25).
    WithOffset(0)

sql, args := sqlquery.FindAllQuery("tasks", options)
// SELECT * FROM tasks WHERE priority >= $1 AND status = $2 
// ORDER BY priority DESC, created_at ASC LIMIT 25 OFFSET 0
```

## Row Locking

Support for FOR UPDATE with various locking modes (database-specific):

```go
// Basic FOR UPDATE
options := sqlquery.NewFindOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("id", 1).
    WithForUpdate("")

sql, args := sqlquery.FindQuery("accounts", options)
// SELECT * FROM accounts WHERE id = $1 FOR UPDATE

// FOR UPDATE NOWAIT (PostgreSQL/MySQL)
options := sqlquery.NewFindOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("id", 1).
    WithForUpdate("NOWAIT")

sql, args := sqlquery.FindQuery("accounts", options)
// SELECT * FROM accounts WHERE id = $1 FOR UPDATE NOWAIT

// FOR UPDATE SKIP LOCKED (PostgreSQL/MySQL 8.0+)
options := sqlquery.NewFindAllOptions(sqlquery.MySQLFlavor).
    WithFilter("status", "pending").
    WithLimit(10).
    WithForUpdate("SKIP LOCKED")

sql, args := sqlquery.FindAllQuery("queue_items", options)
// SELECT * FROM queue_items WHERE status = ? LIMIT 10 FOR UPDATE SKIP LOCKED

// Useful for job queues
options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("status", "pending").
    WithOrderBy("priority DESC, created_at ASC").
    WithLimit(1).
    WithForUpdate("SKIP LOCKED")

sql, args := sqlquery.FindAllQuery("jobs", options)
// SELECT * FROM jobs WHERE status = $1 ORDER BY priority DESC, created_at ASC 
// LIMIT 1 FOR UPDATE SKIP LOCKED
```

## Struct-Based Queries

Use struct tags to control which fields are included in INSERT and UPDATE operations:

```go
type User struct {
    ID        int       `db:"id" insert:"insert"`
    Name      string    `db:"name" insert:"insert" update:"update"`
    Email     string    `db:"email" insert:"insert" update:"update"`
    Password  string    `db:"password" insert:"insert"`
    CreatedAt time.Time `db:"created_at" insert:"insert"`
    UpdatedAt time.Time `db:"updated_at" update:"update"`
}

// INSERT with "insert" tag
user := User{
    ID:        1,
    Name:      "Alice Smith",
    Email:     "alice@example.com",
    Password:  "hashed_password",
    CreatedAt: time.Now(),
}

sql, args := sqlquery.InsertQuery(
    sqlquery.PostgreSQLFlavor,
    "insert",
    "users",
    &user,
)
// SQL: INSERT INTO users (id, name, email, password, created_at) 
//      VALUES ($1, $2, $3, $4, $5)
// Args: [1, "Alice Smith", "alice@example.com", "hashed_password", <timestamp>]

// UPDATE with "update" tag (only updates name, email, updated_at)
user.Name = "Alice Johnson"
user.Email = "alice.johnson@example.com"
user.UpdatedAt = time.Now()

sql, args := sqlquery.UpdateQuery(
    sqlquery.PostgreSQLFlavor,
    "update",
    "users",
    user.ID,
    &user,
)
// SQL: UPDATE users SET email = $1, name = $2, updated_at = $3 WHERE id = $4
// Args: ["alice.johnson@example.com", "Alice Johnson", <timestamp>, 1]
```

## Database Flavors

sqlquery supports three database flavors with appropriate SQL dialect handling:

```go
// MySQL / MariaDB
sqlquery.MySQLFlavor
// Uses: ? for placeholders, MySQL-specific syntax

// PostgreSQL
sqlquery.PostgreSQLFlavor
// Uses: $1, $2, $3 for placeholders, PostgreSQL-specific syntax

// SQLite
sqlquery.SQLiteFlavor
// Uses: ? for placeholders, SQLite-specific syntax
```

### Flavor Examples

```go
// MySQL
options := sqlquery.NewFindOptions(sqlquery.MySQLFlavor).
    WithFilter("id", 1)
sql, args := sqlquery.FindQuery("users", options)
// SELECT * FROM users WHERE id = ?

// PostgreSQL
options := sqlquery.NewFindOptions(sqlquery.PostgreSQLFlavor).
    WithFilter("id", 1)
sql, args := sqlquery.FindQuery("users", options)
// SELECT * FROM users WHERE id = $1

// SQLite
options := sqlquery.NewFindOptions(sqlquery.SQLiteFlavor).
    WithFilter("id", 1)
sql, args := sqlquery.FindQuery("users", options)
// SELECT * FROM users WHERE id = ?
```

## Complete Examples

### Example 1: User Management System

```go
package main

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/allisson/sqlquery"
    _ "github.com/lib/pq"
)

type User struct {
    ID        int       `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Status    string    `db:"status"`
    CreatedAt time.Time `db:"created_at"`
}

func main() {
    db, _ := sql.Open("postgres", "postgresql://user:pass@localhost/mydb?sslmode=disable")
    defer db.Close()

    // Find active users with pagination
    options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor).
        WithFilter("status", "active").
        WithFilter("created_at.gte", time.Now().AddDate(0, -1, 0)).
        WithOrderBy("created_at DESC").
        WithLimit(20).
        WithOffset(0)

    sql, args := sqlquery.FindAllQuery("users", options)
    rows, _ := db.Query(sql, args...)
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.CreatedAt)
        users = append(users, u)
    }

    fmt.Printf("Found %d active users\n", len(users))
}
```

### Example 2: E-commerce Order Processing

```go
package main

import (
    "database/sql"
    "time"

    "github.com/allisson/sqlquery"
    _ "github.com/go-sql-driver/mysql"
)

func processOrders(db *sql.DB) error {
    // Start transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Lock and fetch pending orders
    options := sqlquery.NewFindAllOptions(sqlquery.MySQLFlavor).
        WithFilter("status", "pending").
        WithFilter("created_at.lt", time.Now().Add(-5*time.Minute)).
        WithOrderBy("priority DESC, created_at ASC").
        WithLimit(10).
        WithForUpdate("SKIP LOCKED")

    sql, args := sqlquery.FindAllQuery("orders", options)
    rows, err := tx.Query(sql, args...)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Process each order
    for rows.Next() {
        var orderID int
        rows.Scan(&orderID)

        // Update order status
        updateOpts := sqlquery.NewUpdateOptions(sqlquery.MySQLFlavor).
            WithAssignment("status", "processing").
            WithAssignment("processed_at", time.Now()).
            WithFilter("id", orderID)

        updateSQL, updateArgs := sqlquery.UpdateWithOptionsQuery("orders", updateOpts)
        _, err = tx.Exec(updateSQL, updateArgs...)
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### Example 3: Data Cleanup Job

```go
package main

import (
    "database/sql"
    "log"
    "time"

    "github.com/allisson/sqlquery"
    _ "github.com/mattn/go-sqlite3"
)

func cleanupOldRecords(db *sql.DB) error {
    // Delete inactive users older than 2 years
    deleteOpts := sqlquery.NewDeleteOptions(sqlquery.SQLiteFlavor).
        WithFilter("status", "inactive").
        WithFilter("last_login.lt", time.Now().AddDate(-2, 0, 0))

    sql, args := sqlquery.DeleteWithOptionsQuery("users", deleteOpts)
    result, err := db.Exec(sql, args...)
    if err != nil {
        return err
    }

    deleted, _ := result.RowsAffected()
    log.Printf("Deleted %d inactive users\n", deleted)

    // Archive completed orders older than 1 year
    archiveOpts := sqlquery.NewFindAllOptions(sqlquery.SQLiteFlavor).
        WithFilter("status", "completed").
        WithFilter("completed_at.lt", time.Now().AddDate(-1, 0, 0)).
        WithLimit(1000)

    sql, args = sqlquery.FindAllQuery("orders", archiveOpts)
    rows, err := db.Query(sql, args...)
    if err != nil {
        return err
    }
    defer rows.Close()

    // Process archival...

    return nil
}
```

### Example 4: Search with Multiple Filters

```go
package main

import (
    "database/sql"
    "fmt"

    "github.com/allisson/sqlquery"
)

type SearchParams struct {
    Query    string
    Category string
    MinPrice float64
    MaxPrice float64
    InStock  bool
    Page     int
    PageSize int
}

func searchProducts(db *sql.DB, params SearchParams) ([]Product, error) {
    options := sqlquery.NewFindAllOptions(sqlquery.PostgreSQLFlavor)

    // Add filters based on search parameters
    if params.Query != "" {
        options = options.WithFilter("name.like", "%"+params.Query+"%")
    }
    if params.Category != "" {
        options = options.WithFilter("category", params.Category)
    }
    if params.MinPrice > 0 {
        options = options.WithFilter("price.gte", params.MinPrice)
    }
    if params.MaxPrice > 0 {
        options = options.WithFilter("price.lte", params.MaxPrice)
    }
    if params.InStock {
        options = options.WithFilter("stock.gt", 0)
    }

    // Add pagination
    offset := (params.Page - 1) * params.PageSize
    options = options.
        WithLimit(params.PageSize).
        WithOffset(offset).
        WithOrderBy("popularity DESC, name ASC")

    sql, args := sqlquery.FindAllQuery("products", options)
    rows, err := db.Query(sql, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        // Scan into product...
        products = append(products, p)
    }

    return products, nil
}
```

## API Reference

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/allisson/sqlquery).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
