package sqlb

import (
	"fmt"
	"strings"
)

type UpdateQuery struct {
	fields  []string
	values  []any
	returns []string
	stmt    string
	clause  string
}

func Update(table string) *UpdateQuery {
	return &UpdateQuery{
		stmt: "update " + table + " set ",
	}
}

func (q *UpdateQuery) Set(field string, value any) *UpdateQuery {
	q.values = append(q.values, value)
	q.fields = append(q.fields, field)
	return q
}

func (q *UpdateQuery) Where(column string, op string, value any) *UpdateQuery {
	if strings.Count(q.clause, "where") > 0 {
		q.clause += " and"
	} else {
		q.clause += " where"
	}
	q.clause += fmt.Sprintf(" %s %s $%d", column, op, len(q.values)+1)
	q.values = append(q.values, value)
	return q
}

func (q *UpdateQuery) Return(fields ...string) *UpdateQuery {
	q.returns = fields
	return q
}

func (q *UpdateQuery) Values() []any {
	return q.values
}

func (q *UpdateQuery) Build() *UpdateQuery {
	for k := range q.fields {
		funcValue, ok := q.values[k].(ValueFunc)
		if ok {
			stmt, value := recFuncValue(funcValue, k+1)
			q.stmt += q.fields[k] + "=" + stmt
			q.values[k] = value
		} else {
			q.stmt += fmt.Sprintf("%s=$%d", q.fields[k], k+1)
		}
		if k < len(q.fields)-1 {
			q.stmt += ","
		}
	}
	if q.clause != "" {
		q.stmt += " " + q.clause
	}
	if len(q.returns) > 0 {
		q.stmt += fmt.Sprintf(" returning %s", strings.Join(q.returns, ","))
	}
	return q
}

func (q *UpdateQuery) Stmt() string {
	return q.stmt
}

func (q *UpdateQuery) Debug() string {
	return Debug(q.stmt, q.values...)
}
