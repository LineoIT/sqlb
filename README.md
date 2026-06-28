# sqlb

A fluent SQL query builder for Go targeting PostgreSQL. All user-supplied values are bound as positional parameters (`$1`, `$2`, ...) — the builder never interpolates values directly into the query string, which eliminates SQL injection at the value level. Table names, column names, and identifiers are validated by the library to prevent injection via structural elements.

## Installation

```bash
go get -u github.com/LineoIT/sqlb
```

Requires Go 1.19 or later.

---

## Table of contents

- [SELECT](#select)
- [INSERT](#insert)
- [Batch INSERT](#batch-insert)
- [UPDATE](#update)
- [DELETE](#delete)
- [WHERE conditions](#where-conditions)
- [Expressions](#expressions)
- [Scoped conditions](#scoped-conditions)
- [Raw fragments](#raw-fragments)
- [Subqueries](#subqueries)
- [JOINs](#joins)
- [GROUP BY and HAVING](#group-by-and-having)
- [Common Table Expressions (WITH)](#common-table-expressions-with)
- [Value functions](#value-functions)
- [Text search helpers](#text-search-helpers)
- [Using the built query](#using-the-built-query)
- [Debug and utilities](#debug-and-utilities)
- [API reference](#api-reference)

---

## SELECT

```go
q := sqlb.Select("id", "name", "email").
    From("users").
    Where("active", true).
    OrderBy("created_at").
    Sort(sqlb.DESC).
    Limit(20).
    Offset(0).
    Build()

fmt.Println(q.Stmt())
// select id,name,email from users where active=$1 order by created_at DESC limit 20;

fmt.Println(q.Args())
// [true]
```

Pass no columns to `Select` to get `SELECT *`:

```go
q := sqlb.Select().From("users").Build()
// select * from users;
```

Add `DISTINCT`:

```go
q := sqlb.Select("email").Distinct().From("users").Build()
// select distinct email from users;
```

Select from multiple tables:

```go
q := sqlb.Select("p.*", "a.city").From("persons p", "addresses a").Build()
// select p.*,a.city from persons p,addresses a;
```

### Starting from a raw SQL base

Use `SQL()` when the base query is complex and cannot be expressed with `Select().From()` alone (e.g. queries with subselects in the column list):

```go
base := `SELECT notifications.id, notifications.subject,
         COALESCE(json_agg(likes.*) FILTER (WHERE likes.pivot_id IS NOT NULL), '[]'::json) AS likes
         FROM notifications`

q := sqlb.SQL(base).
    Where("notifications.user_id", userID).
    Where(sqlb.Greater("notifications.created_at", since)).
    Build()
```

---

## INSERT

```go
q := sqlb.Insert("users").
    Value("email", "alice@example.com").
    Value("age", 30).
    Value("active", true).
    Return("id", "created_at").
    Build()

fmt.Println(q.Stmt())
// insert into users (email,age,active) values($1,$2,$3) returning id,created_at

db.QueryRow(q.Stmt(), q.Values()...)
```

### UPSERT (ON CONFLICT)

Do nothing on conflict:

```go
q := sqlb.Insert("products").
    Value("sku", "ABC-123").
    Value("price", 9.99).
    OnConflict("sku").Nothing().
    Build()
// insert into products (sku,price) values($1,$2) on conflict(sku) do nothing
```

Update on conflict:

```go
q := sqlb.Insert("users").
    Value("email", "alice@example.com").
    Value("age", 30).
    OnConflict("email").Update().
    Set("age", 31).
    Return("id").
    Build()
// insert into users (email,age) values($1,$2) on conflict(email) do update set age=$3 returning id
```

---

## Batch INSERT

Insert multiple rows in a single statement:

```go
q := sqlb.BatchInsert("users", "email", "age").
    Values("alice@example.com", 30).
    Values("bob@example.com", 25).
    Values("carol@example.com", 28).
    Return("id").
    Build()

fmt.Println(q.Stmt())
// insert into users(email,age) values($1, $2),($3, $4),($5, $6) returning id

db.Query(q.Stmt(), q.Args()...)
```

Columns can also be declared separately with `Columns()`:

```go
q := sqlb.BatchInsert("users").
    Columns("email", "age").
    Values("alice@example.com", 30).
    Build()
```

---

## UPDATE

```go
q := sqlb.Update("users").
    Set("name", "Alice").
    Set("age", 31).
    Where("id", 1).
    Return("updated_at").
    Build()

fmt.Println(q.Stmt())
// update users set name=$1,age=$2 where id=$3 returning updated_at;
```

---

## DELETE

```go
q := sqlb.Delete("users").
    Where("id", 5).
    Build()

fmt.Println(q.Stmt())
// delete from users where id=$1
```

With multiple conditions:

```go
q := sqlb.Delete("sessions").
    Where("user_id", userID).
    Where(sqlb.Less("expires_at", time.Now())).
    Build()
```

---

## WHERE conditions

Chain `Where()` calls to add AND conditions, and `Or()` to add OR conditions. `Or()` must follow a `Where()` or `Having()` call.

```go
q := sqlb.Select("*").From("orders").
    Where("status", "paid").
    Where(sqlb.Greater("total", 100)).
    Or(sqlb.IsNull("deleted_at")).
    Build()
// select * from orders where status=$1 and total > $2 or deleted_at is null;
```

### IN and NOT IN

```go
q := sqlb.Select("*").From("users").
    Where(sqlb.In("role", []string{"admin", "moderator"})).
    Build()
// select * from users where role in ($1,$2);
```

```go
q := sqlb.Select("*").From("users").
    Where(sqlb.NotIn("status", []string{"banned", "suspended"})).
    Build()
// select * from users where status not in ($1,$2);
```

### BETWEEN

```go
q := sqlb.Select("*").From("orders").
    Where(sqlb.Between("total", []float64{100.0, 500.0})).
    Build()
// select * from orders where total between $1 and $2;
```

### NULL checks

```go
q := sqlb.Select("*").From("users").
    Where(sqlb.IsNull("deleted_at")).
    Build()
// select * from users where deleted_at is null;

q = sqlb.Select("*").From("users").
    Where(sqlb.IsNotNull("email")).
    Build()
// select * from users where email is not null;
```

### LIKE and ILIKE

```go
// Case-sensitive
q := sqlb.Select("*").From("products").
    Where(sqlb.Like("sku", "ABC%")).
    Build()
// select * from products where sku like $1;

// Case-insensitive
q = sqlb.Select("*").From("users").
    Where(sqlb.Ilike("name", "%alice%")).
    Build()
// select * from users where name ilike $1;
```

### Comparison operators

```go
sqlb.Equal("status", "active")        // status = $N
sqlb.NotEqual("status", "banned")     // status <> $N
sqlb.Greater("age", 18)               // age > $N
sqlb.GreaterOrEqual("score", 90)      // score >= $N
sqlb.Less("stock", 10)                // stock < $N
sqlb.LessOrEqual("price", 999.99)     // price <= $N
```

```go
q := sqlb.Select("*").From("products").
    Where(sqlb.GreaterOrEqual("price", 10.0)).
    Where(sqlb.Less("price", 100.0)).
    Build()
// select * from products where price >= $1 and price < $2;
```

### IS and IS NOT

```go
sqlb.Is("verified", nil)       // verified is null
sqlb.Is("verified", true)      // verified is $N
sqlb.IsNot("role", nil)        // role is not null
sqlb.IsNot("role", "guest")    // role is not $N
```

### Generic Expression

`Expression` maps a string operator name to the correct expression type:

```go
sqlb.Expression("role", "in", []string{"admin", "user"})
sqlb.Expression("age", "between", []int{18, 65})
sqlb.Expression("email", "ilike", "%@company.com")
sqlb.Expression("deleted_at", "is null", nil)
```

---

## Scoped conditions

Use `Scope()` and `EndScope()` to wrap conditions in parentheses:

```go
q := sqlb.Select("*").From("users").
    Where("active", true).
    Scope().
        Or("role", "admin").
        Or("role", "moderator").
    EndScope().
    Build()
// select * from users where active=$1 or (role=$2 or role=$3);
```

---

## Raw fragments

When the builder's built-in methods are not sufficient, inject raw SQL with `WhereRaw`, `OrRaw`, or `Raw`. Placeholders in raw strings must start at `$1`; the builder renumbers them automatically to avoid collisions with previously bound arguments.

```go
q := sqlb.Select("*").From("users").
    Where("active", true).
    WhereRaw("name ilike $1", "%john%").
    Build()
// select * from users where active=$1 and name ilike $2;
// args: [true, "%john%"]
```

```go
q := sqlb.Select("*").From("users").
    Where("role", "admin").
    OrRaw("email ilike $1", "%@company.com").
    Build()
// select * from users where role=$1 or email ilike $2;
```

`Raw` appends an arbitrary SQL fragment without any WHERE/AND/OR wrapping:

```go
q := sqlb.SQL("select * from stops").
    Raw("where latitude is not null").
    Build()
```

`OrderByRaw` accepts a raw expression for ORDER BY:

```go
q := sqlb.Select("*").From("users").
    OrderByRaw("lower(name)").
    Build()
// select * from users order by lower(name);
```

---

## Subqueries

### WHERE EXISTS / NOT EXISTS

```go
sub := sqlb.Select("1").From("orders").
    Where("orders.user_id", userID).
    Where("orders.status", "paid").
    Build()

q := sqlb.Select("*").From("users").
    WhereExists(sub).
    Build()
// select * from users where exists (select 1 from orders where orders.user_id=$1 and orders.status=$2);
```

```go
q := sqlb.Select("*").From("users").
    WhereNotExists(sub).
    Build()
// select * from users where not exists (...);
```

### WHERE column IN (subquery)

```go
sub := sqlb.Select("user_id").From("orders").
    Where(sqlb.Greater("total", 1000)).
    Build()

q := sqlb.Select("*").From("users").
    WhereInSubquery("id", sub).
    Build()
// select * from users where id in (select user_id from orders where total > $1);
```

```go
q := sqlb.Select("*").From("users").
    WhereNotInSubquery("id", sub).
    Build()
// select * from users where id not in (...);
```

---

## JOINs

```go
q := sqlb.Select("u.id", "u.name", "o.total").
    From("users u").
    InnerJoin("orders o").On("u.id", "o.user_id").
    Where("u.active", true).
    Build()
// select u.id,u.name,o.total from users u inner join orders o on u.id=o.user_id where u.active=$1;
```

Available join methods:

| Method | SQL |
|---|---|
| `Join(table)` | `JOIN table` |
| `InnerJoin(table)` | `INNER JOIN table` |
| `LeftJoin(table)` | `LEFT JOIN table` |
| `RightJoin(table)` | `RIGHT JOIN table` |
| `FullJoin(table)` | `FULL JOIN table` |

For join conditions that cannot be expressed as `col = col`, use `Raw`:

```go
q := sqlb.Select("*").From("events e").
    LeftJoin("categories c").
    Raw("on e.category_id = c.id and c.active = true").
    Where("e.published", true).
    Build()
```

---

## GROUP BY and HAVING

```go
q := sqlb.Select("role", "count(*)").From("users").
    GroupBy("role").
    Having(sqlb.Greater("count(*)", 5)).
    Build()
// select role,count(*) from users group by role having count(*) > $1;
```

Multiple GROUP BY columns:

```go
q := sqlb.Select("department", "role", "count(*)").From("employees").
    GroupBy("department", "role").
    Having(sqlb.GreaterOrEqual("count(*)", 3)).
    Build()
```

---

## Common Table Expressions (WITH)

### Single CTE

```go
active := sqlb.Select("id", "name").From("users").Where("active", true).Build()

q := sqlb.With("active_users", active).
    Query(sqlb.Select("*").From("active_users").Build()).
    Build()
// with active_users as (select id,name from users where active=$1) select * from active_users;
```

### Multiple CTEs

```go
recentOrders := sqlb.Select("id", "user_id").From("orders").Where("status", "paid").Build()
activeUsers  := sqlb.Select("id", "name").From("users").Where("active", true).Build()

q := sqlb.NewCTE().
    Add("recent_orders", recentOrders).
    Add("active_users", activeUsers).
    Query(
        sqlb.Select("u.name", "o.id").
            From("active_users u").
            InnerJoin("recent_orders o").On("u.id", "o.user_id").
            Build(),
    ).
    Build()

// with recent_orders as (...), active_users as (...) select u.name,o.id from ...
// args: ["paid", true]  — placeholders are renumbered automatically
```

### Recursive CTE

```go
q := sqlb.WithRecursive("nums", sqlb.SQL("select 1 as n")).
    AddRaw("", "select n+1 from nums where n < 10").
    Query(sqlb.Select("*").From("nums").Build()).
    Build()
// with recursive nums as (select 1 as n), as (select n+1 from nums where n < 10) select * from nums;
```

---

## Value functions

`Coalesce` and `Nullif` can be used as values in `Insert.Value()` or `Update.Set()`:

```go
q := sqlb.Insert("users").
    Value("name", sqlb.Coalesce(inputName, "anonymous")).
    Value("score", sqlb.Nullif(score, 0)).
    Build()
// insert into users (name,score) values(coalesce($1,anonymous),nullif($2,0))
```

With a type cast:

```go
sqlb.Coalesce(value, "default", "text")  // coalesce($N::text, default)
```

---

## Text search helpers

These helpers are shorthand for `ILIKE` patterns. They continue whichever clause is currently active (`WHERE` or `HAVING`).

```go
q := sqlb.Select("*").From("users").
    Where("active", true).
    Contains("name", "alice").
    Build()
// select * from users where active=$1 and name ilike $2;
// args: [true, "%alice%"]
```

```go
q := sqlb.Select("*").From("users").
    Where("active", true).
    StartWith("email", "admin").
    Build()
// args: [true, "admin%"]
```

```go
q := sqlb.Select("*").From("users").
    Where("active", true).
    EndWith("email", "@acme.com").
    Build()
// args: [true, "%@acme.com"]
```

---

## Using the built query

Every builder exposes `Stmt()` and `Args()` (or `Values()` for INSERT/UPDATE), which map directly to the standard `database/sql` interface.

```go
import "database/sql"

q := sqlb.Select("id", "name").From("users").Where("active", true).Build()

rows, err := db.Query(q.Stmt(), q.Args()...)
```

For `pgx`:

```go
rows, err := pool.Query(ctx, q.Stmt(), q.Args()...)
```

Check for builder errors before executing:

```go
q := sqlb.Select("*").From("users").OrderBy(userInput)
if err := q.Error(); err != nil {
    return err
}
result := q.Build()
```

---

## Debug and utilities

### Debug()

Returns a human-readable version of the query with placeholders substituted by their actual values. For development and logging only — never pass this string to a database driver.

```go
q := sqlb.Select("*").From("orders").
    Where("status", "paid").
    Where(sqlb.Greater("total", 500)).
    Build()

fmt.Println(q.Debug())
// select * from orders where status=paid and total > 500;
```

### CleanSQL()

Normalises whitespace in a multi-line SQL string (strips tabs, collapses newlines, trims leading/trailing spaces):

```go
cleaned := sqlb.CleanSQL(`
    SELECT id, name
    FROM users
    WHERE active = true
`)
// "SELECT id, name FROM users WHERE active = true"
```

---

## API reference

### Entry points

| Function | Returns | Description |
|---|---|---|
| `Select(cols ...string)` | `*selectQuery` | Start a SELECT query |
| `SQL(base string)` | `*QueryBuilder` | Start from a raw SQL base |
| `Insert(table string)` | `*InsertQuery` | Start an INSERT query |
| `BatchInsert(table string, cols ...string)` | `*BatchInsertQuery` | Start a multi-row INSERT |
| `Update(table string)` | `*UpdateQuery` | Start an UPDATE query |
| `Delete(table string)` | `*DeleteQuery` | Start a DELETE query |
| `NewCTE()` | `*CTEBuilder` | Start a CTE (WITH) builder |
| `With(name, q)` | `*CTEBuilder` | Shorthand for `NewCTE().Add(name, q)` |
| `WithRecursive(name, q)` | `*CTEBuilder` | Shorthand for recursive CTE |

### QueryBuilder methods

| Method | Description |
|---|---|
| `Where(col, val)` | Add AND WHERE condition |
| `Or(col, val)` | Add OR condition (after Where or Having) |
| `WhereRaw(raw, args...)` | Raw AND WHERE fragment |
| `OrRaw(raw, args...)` | Raw OR fragment |
| `WhereExists(sub)` | WHERE EXISTS (subquery) |
| `WhereNotExists(sub)` | WHERE NOT EXISTS (subquery) |
| `WhereInSubquery(col, sub)` | WHERE col IN (subquery) |
| `WhereNotInSubquery(col, sub)` | WHERE col NOT IN (subquery) |
| `Scope()` / `EndScope()` | Open / close a parenthesised group |
| `GroupBy(cols...)` | GROUP BY clause |
| `Having(col, val)` | HAVING condition |
| `OrderBy(cols...)` | ORDER BY (validated identifiers) |
| `OrderByRaw(raw)` | ORDER BY with a raw expression |
| `Sort(dir)` | Add ASC or DESC to ORDER BY |
| `Limit(n)` | LIMIT clause |
| `Offset(n)` | OFFSET clause |
| `Take(limit, offset)` | Set both LIMIT and OFFSET |
| `Join(table)` | JOIN |
| `InnerJoin(table)` | INNER JOIN |
| `LeftJoin(table)` | LEFT JOIN |
| `RightJoin(table)` | RIGHT JOIN |
| `FullJoin(table)` | FULL JOIN |
| `On(fk, ref)` | ON clause for the preceding join |
| `Raw(fragment)` | Append raw SQL |
| `Contains(col, val)` | WHERE col ILIKE %val% |
| `StartWith(col, val)` | WHERE col ILIKE val% |
| `EndWith(col, val)` | WHERE col ILIKE %val |
| `Build()` | Finalise and return the builder |
| `Stmt()` | Return the built SQL string |
| `Args()` | Return the bound argument slice |
| `Debug()` | Return SQL with values inlined (for logging) |
| `Error()` | Return any accumulated error |

### Expression functions

| Function | SQL equivalent |
|---|---|
| `Equal(col, val)` | `col = $N` |
| `NotEqual(col, val)` | `col <> $N` |
| `Greater(col, val)` | `col > $N` |
| `GreaterOrEqual(col, val)` | `col >= $N` |
| `Less(col, val)` | `col < $N` |
| `LessOrEqual(col, val)` | `col <= $N` |
| `In(col, vals)` | `col in ($1,$2,...)` |
| `NotIn(col, vals)` | `col not in ($1,$2,...)` |
| `Between(col, [from,to])` | `col between $N and $M` |
| `IsNull(col)` | `col is null` |
| `IsNotNull(col)` | `col is not null` |
| `Is(col, val)` | `col is $N` / `col is null` |
| `IsNot(col, val)` | `col is not $N` / `col is not null` |
| `Like(col, val)` | `col like $N` |
| `Ilike(col, val)` | `col ilike $N` |
| `Expression(col, op, val)` | generic dispatcher |
| `Coalesce(val, alt, cast?)` | `coalesce($N, alt)` |
| `Nullif(val, alt, cast?)` | `nullif($N, alt)` |

---

## License

See [LICENCE](LICENCE).
