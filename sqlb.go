package sqlb

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	limit, offset int64
	orderBy, sort string
	stmt          string
	args          []interface{}
	currentTag    string
	error         error
}

func SQL(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		stmt: baseQuery,
		args: make([]interface{}, 0),
	}
}

func (q *QueryBuilder) Build() *QueryBuilder {
	q.stmt = strings.ReplaceAll(q.stmt, beginScope, "(")
	q.stmt = strings.ReplaceAll(q.stmt, endScope, ")")
	if q.orderBy != "" {
		q.stmt += fmt.Sprintf(" order by %s", q.orderBy)
		if strings.ToUpper(q.sort) == DESC {
			q.stmt += " " + q.sort
		}
	}
	if q.limit > 0 {
		q.stmt += fmt.Sprintf(` limit %d`, q.limit)
	}
	if q.offset > 0 {
		q.stmt += fmt.Sprintf(` offset %d`, q.offset)
	}
	q.stmt += ";"
	return q
}

func (q *QueryBuilder) Args() []interface{} {
	return q.args
}

func (q *QueryBuilder) Stmt() string {
	return q.stmt
}

func (q *QueryBuilder) Where(column string, value interface{}) *QueryBuilder {
	if q.error == nil {
		if err := validateFilterColumn(column); err != nil {
			q.error = err
			return q
		}
	}
	q.filter(whereVar, column, value)
	return q
}

func (q *QueryBuilder) Take(limit, offset int64) *QueryBuilder {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *QueryBuilder) Limit(limit int64) *QueryBuilder {
	q.limit = limit
	return q
}

func (q *QueryBuilder) Offset(offset int64) *QueryBuilder {
	q.offset = offset
	return q
}

// OrderBy adds an ORDER BY clause. Each column must be a plain identifier (table.column).
// For expressions like LOWER(name) use OrderByRaw.
func (q *QueryBuilder) OrderBy(columns ...string) *QueryBuilder {
	for _, col := range columns {
		if err := validateIdentifier(col); err != nil {
			q.error = err
			return q
		}
	}
	q.orderBy = strings.Join(columns, ",")
	return q
}

// OrderByRaw adds an ORDER BY clause with a raw expression. The caller is responsible
// for ensuring the expression is safe.
func (q *QueryBuilder) OrderByRaw(raw string) *QueryBuilder {
	q.orderBy = raw
	return q
}

func (q *QueryBuilder) Sort(sort string) *QueryBuilder {
	q.sort = sort
	return q
}

func (q *QueryBuilder) Or(column string, value interface{}) *QueryBuilder {
	if q.currentTag == whereVar || q.currentTag == havingVar {
		if q.error == nil {
			if err := validateFilterColumn(column); err != nil {
				q.error = err
				return q
			}
		}
		q.normalizeFilter("or", column, value)
		return q
	}
	q.error = fmt.Errorf("sqlb: Or must be called after Where or Having")
	return q
}

// GroupBy adds a GROUP BY clause. Each column must be a plain identifier.
func (q *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	for _, col := range columns {
		if err := validateIdentifier(col); err != nil {
			q.error = err
			return q
		}
	}
	cols := strings.Join(columns, ",")
	if isClauseExist(q.stmt, groupByVar) {
		q.stmt += "," + cols
	} else {
		q.stmt += " " + groupByVar + " " + cols
	}
	q.currentTag = groupByVar
	return q
}

func (q *QueryBuilder) Having(column string, value interface{}) *QueryBuilder {
	if q.currentTag == havingVar || q.currentTag == groupByVar || isClauseExist(q.stmt, groupByVar) {
		if q.error == nil {
			if err := validateFilterColumn(column); err != nil {
				q.error = err
				return q
			}
		}
		q.filter(havingVar, column, value)
		return q
	}
	q.error = fmt.Errorf("sqlb: Having must be called after GroupBy")
	return q
}

// Raw appends a raw SQL fragment. Values must already be parameterised by the caller.
func (q *QueryBuilder) Raw(raw string) *QueryBuilder {
	q.stmt += " " + raw
	return q
}

func (q *QueryBuilder) Scope() *QueryBuilder {
	q.stmt += beginScope
	return q
}

func (q *QueryBuilder) EndScope() *QueryBuilder {
	q.stmt += endScope
	return q
}

// WhereRaw appends a raw WHERE condition. Placeholders in raw must start at $1; they
// are renumbered automatically to follow any previously bound arguments.
func (q *QueryBuilder) WhereRaw(raw string, args ...interface{}) *QueryBuilder {
	if len(args) > 0 {
		raw = renumberArgs(raw, len(q.args))
		q.args = append(q.args, args...)
	}
	query := strings.ToLower(q.stmt)
	if strings.LastIndex(query, "where") > strings.LastIndex(query, "from") {
		if strings.Contains(q.stmt, beginScope) {
			q.stmt = strings.ReplaceAll(q.stmt, beginScope, "")
			q.stmt += " and ("
		} else {
			q.stmt += " and"
		}
	} else {
		q.stmt += " where"
		q.currentTag = whereVar
	}
	q.stmt += " " + raw
	return q
}

// OrRaw appends a raw OR condition. Placeholders in raw must start at $1; they are
// renumbered automatically.
func (q *QueryBuilder) OrRaw(raw string, args ...interface{}) *QueryBuilder {
	if len(args) > 0 {
		raw = renumberArgs(raw, len(q.args))
		q.args = append(q.args, args...)
	}
	query := strings.ToLower(q.stmt)
	if strings.LastIndex(query, "where") > strings.LastIndex(query, "from") {
		if strings.Contains(q.stmt, beginScope) {
			q.stmt = strings.ReplaceAll(q.stmt, beginScope, "")
			q.stmt += " or ("
		} else {
			q.stmt += " or"
		}
	} else {
		q.stmt += " where"
		q.currentTag = whereVar
	}
	q.stmt += " " + raw
	return q
}

// WhereExists adds a WHERE EXISTS (subquery) clause.
func (q *QueryBuilder) WhereExists(sub *QueryBuilder) *QueryBuilder {
	subStmt := strings.TrimSuffix(strings.TrimSpace(sub.Stmt()), ";")
	subStmt = renumberArgs(subStmt, len(q.args))
	q.args = append(q.args, sub.Args()...)
	return q.WhereRaw("exists ("+subStmt+")")
}

// WhereNotExists adds a WHERE NOT EXISTS (subquery) clause.
func (q *QueryBuilder) WhereNotExists(sub *QueryBuilder) *QueryBuilder {
	subStmt := strings.TrimSuffix(strings.TrimSpace(sub.Stmt()), ";")
	subStmt = renumberArgs(subStmt, len(q.args))
	q.args = append(q.args, sub.Args()...)
	return q.WhereRaw("not exists ("+subStmt+")")
}

// WhereInSubquery adds WHERE column IN (subquery).
func (q *QueryBuilder) WhereInSubquery(column string, sub *QueryBuilder) *QueryBuilder {
	if err := validateIdentifier(column); err != nil {
		q.error = err
		return q
	}
	subStmt := strings.TrimSuffix(strings.TrimSpace(sub.Stmt()), ";")
	subStmt = renumberArgs(subStmt, len(q.args))
	q.args = append(q.args, sub.Args()...)
	return q.WhereRaw(column+" in ("+subStmt+")")
}

// WhereNotInSubquery adds WHERE column NOT IN (subquery).
func (q *QueryBuilder) WhereNotInSubquery(column string, sub *QueryBuilder) *QueryBuilder {
	if err := validateIdentifier(column); err != nil {
		q.error = err
		return q
	}
	subStmt := strings.TrimSuffix(strings.TrimSpace(sub.Stmt()), ";")
	subStmt = renumberArgs(subStmt, len(q.args))
	q.args = append(q.args, sub.Args()...)
	return q.WhereRaw(column+" not in ("+subStmt+")")
}

func (q *QueryBuilder) Contains(col string, value string) *QueryBuilder {
	c, v := Ilike(col, "%"+value+"%")
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) StartWith(col string, value string) *QueryBuilder {
	c, v := Ilike(col, value+"%")
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) EndWith(col string, value string) *QueryBuilder {
	c, v := Ilike(col, "%"+value)
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) matchAny(col string, value any) {
	if q.currentTag == havingVar {
		q.Having(col, value)
	} else if q.currentTag == whereVar {
		q.Where(col, value)
	}
}

func (q *QueryBuilder) Debug() string {
	return Debug(q.stmt, q.args...)
}

func (q *QueryBuilder) Error() error {
	return q.error
}

func (q *QueryBuilder) filter(clause string, column string, value ...interface{}) {
	queryFilter(&q.error, &q.stmt, &q.currentTag, &q.args, clause, column, value...)
}

func (q *QueryBuilder) normalizeFilter(clauseType string, column string, value ...interface{}) {
	filterNormalizer(&q.error, &q.stmt, &q.args, clauseType, column, value...)
}
