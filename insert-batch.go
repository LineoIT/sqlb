package sqlb

import (
	"fmt"
	"strconv"
	"strings"
)

type BatchInsertQuery struct {
	fields   []string
	values   []any
	returns  []string
	stmt     string
	colsSize int
}

func BatchInsert(table string, cols ...string) *BatchInsertQuery {
	var stmt string
	if len(cols) > 0 {
		stmt = "insert into " + table + "(" + strings.Join(cols, ",") + ")"
	} else {
		stmt = "insert into " + table
	}
	return &BatchInsertQuery{
		stmt: stmt,
	}
}

func (q *BatchInsertQuery) Columns(cols ...string) *BatchInsertQuery {
	if q.colsSize > 0 {
		return q
	}
	q.colsSize = len(cols)
	q.stmt += "(" + strings.Join(cols, ",") + ")"
	return q
}

func (q *BatchInsertQuery) Values(values ...any) *BatchInsertQuery {
	if len(values) == 0 {
		return q
	}
	s := ""
	for i := 1; i <= len(values); i++ {
		s += "$" + strconv.Itoa(len(q.values)+i)
		if i < len(values) {
			s += ", "
		}
	}
	q.fields = append(q.fields, "("+s+")")
	q.values = append(q.values, values...)
	return q
}

func (q *BatchInsertQuery) Return(fields ...string) *BatchInsertQuery {
	q.returns = fields
	return q
}

func (q *BatchInsertQuery) Args() []any {
	return q.values
}

func (q *BatchInsertQuery) Build() *BatchInsertQuery {
	q.stmt += " values" + strings.Join(q.fields, ",")
	if len(q.returns) > 0 {
		q.stmt += fmt.Sprintf(" returning %s", strings.Join(q.returns, ","))
	}
	q.stmt += ";"
	return q
}

func (q *BatchInsertQuery) Stmt() string {
	return q.stmt
}

func (q *BatchInsertQuery) Debug() string {
	return Debug(q.stmt, q.values...)
}
