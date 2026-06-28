package sqlb

import (
	"fmt"
	"strings"
)

// CTEBuilder accumulates named sub-queries for a WITH (Common Table Expression) clause.
//
// Usage:
//
//	active := Select("id", "name").From("users").Where("active", true)
//	q := NewCTE().
//	    Add("active_users", active).
//	    Recursive().
//	    Query(Select("*").From("active_users")).
//	    Build()
type CTEBuilder struct {
	ctes      []cteDef
	recursive bool
	argOffset int
}

type cteDef struct {
	name string
	stmt string
	args []any
}

// NewCTE creates an empty CTEBuilder.
func NewCTE() *CTEBuilder {
	return &CTEBuilder{}
}

// Add appends a named CTE sub-query. The sub-query's arg placeholders are renumbered
// to follow any previously accumulated args.
func (c *CTEBuilder) Add(name string, q interface{ Stmt() string; Args() []interface{} }) *CTEBuilder {
	stmt := strings.TrimSuffix(strings.TrimSpace(q.Stmt()), ";")
	stmt = renumberArgs(stmt, c.argOffset)
	args := q.Args()
	c.ctes = append(c.ctes, cteDef{name: name, stmt: stmt, args: args})
	c.argOffset += len(args)
	return c
}

// AddRaw appends a named CTE defined by a raw SQL string (no args).
func (c *CTEBuilder) AddRaw(name, rawSQL string) *CTEBuilder {
	c.ctes = append(c.ctes, cteDef{name: name, stmt: rawSQL})
	return c
}

// Recursive marks the WITH clause as WITH RECURSIVE.
func (c *CTEBuilder) Recursive() *CTEBuilder {
	c.recursive = true
	return c
}

// Query attaches the main SELECT query that follows the WITH clause and returns a
// *QueryBuilder whose arg list is the concatenation of all CTE args followed by the
// main query args.
func (c *CTEBuilder) Query(q *QueryBuilder) *QueryBuilder {
	return c.build(q.Stmt(), q.Args())
}

// QuerySelect is a convenience wrapper that calls Select(...).From(...) inline.
func (c *CTEBuilder) QuerySelect(cols string, from string) *CTESelectBuilder {
	return &CTESelectBuilder{cte: c, cols: cols, from: from}
}

func (c *CTEBuilder) build(mainStmt string, mainArgs []interface{}) *QueryBuilder {
	mainStmt = strings.TrimSuffix(strings.TrimSpace(mainStmt), ";")
	mainStmt = renumberArgs(mainStmt, c.argOffset)

	parts := make([]string, len(c.ctes))
	var allArgs []any
	for i, def := range c.ctes {
		parts[i] = fmt.Sprintf("%s as (%s)", def.name, def.stmt)
		allArgs = append(allArgs, def.args...)
	}
	allArgs = append(allArgs, mainArgs...)

	keyword := "with"
	if c.recursive {
		keyword = "with recursive"
	}
	return &QueryBuilder{
		stmt: keyword + " " + strings.Join(parts, ", ") + " " + mainStmt,
		args: allArgs,
	}
}

// CTESelectBuilder is a fluent helper returned by QuerySelect.
type CTESelectBuilder struct {
	cte  *CTEBuilder
	cols string
	from string
}

func (b *CTESelectBuilder) Build() *QueryBuilder {
	q := Select(b.cols).From(b.from)
	return b.cte.Query(q)
}

// ---- Convenience top-level functions ----------------------------------------

// With is shorthand for NewCTE().Add(name, q).
func With(name string, q interface{ Stmt() string; Args() []interface{} }) *CTEBuilder {
	return NewCTE().Add(name, q)
}

// WithRecursive is shorthand for NewCTE().Recursive().Add(name, q).
func WithRecursive(name string, q interface{ Stmt() string; Args() []interface{} }) *CTEBuilder {
	return NewCTE().Recursive().Add(name, q)
}
