package sqlb

import (
	"fmt"
	"strings"
)

type InsertQuery struct {
	fields  []string
	values  []any
	returns []string
	stmt    string
}

func Insert(table string) *InsertQuery {
	return &InsertQuery{
		stmt: "insert into " + table,
	}
}

func (q *InsertQuery) Value(field string, value any) *InsertQuery {
	q.values = append(q.values, value)
	q.fields = append(q.fields, field)
	return q
}

func (q *InsertQuery) Return(fields ...string) *InsertQuery {
	q.returns = fields
	return q
}

func (q *InsertQuery) Values() []any {
	return q.values
}

func (q *InsertQuery) Build() *InsertQuery {
	q.stmt += fmt.Sprintf(" (%s) values(", strings.Join(q.fields, ","))
	for k, v := range q.values {
		funcValue, ok := v.(ValueFunc)
		if ok {
			stmt, value := recFuncValue(funcValue, k+1)
			q.stmt += stmt
			q.values[k] = value
		} else {
			q.stmt += fmt.Sprintf("$%d", k+1)
		}
		if k < len(q.values)-1 {
			q.stmt += ","
		}
	}
	q.stmt += ")"
	if len(q.returns) > 0 {
		q.stmt += fmt.Sprintf(" returning %s", strings.Join(q.returns, ","))
	}
	return q
}

func (q *InsertQuery) Stmt() string {
	return q.stmt
}

func (q *InsertQuery) Debug() string {
	return Debug(q.stmt, q.values...)
}
