package sqlb

import (
	"fmt"
	"strings"
)

type UpdateQuery struct {
	fields        []string
	values        []any
	returns       []string
	stmt          string
	whereStmt     string
	currentTag    string
	error         error
	limit, offset int64
	orderBy, sort string
}

func Update(table string) *UpdateQuery {
	q := &UpdateQuery{}
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.stmt = "update " + table + " set "
	return q
}

func (q *UpdateQuery) Set(field string, value any) *UpdateQuery {
	if q.error == nil {
		if err := validateIdentifier(field); err != nil {
			q.error = err
			return q
		}
	}
	q.values = append(q.values, value)
	q.fields = append(q.fields, field)
	return q
}

func (q *UpdateQuery) Where(column string, value interface{}) *UpdateQuery {
	if q.error == nil {
		if err := validateFilterColumn(column); err != nil {
			q.error = err
			return q
		}
	}
	q.clause(whereVar, column, value)
	return q
}

func (q *UpdateQuery) Having(column string, value interface{}) *UpdateQuery {
	if q.currentTag == havingVar || q.currentTag == groupByVar || isClauseExist(q.whereStmt, groupByVar) {
		if q.error == nil {
			if err := validateFilterColumn(column); err != nil {
				q.error = err
				return q
			}
		}
		q.clause(havingVar, column, value)
		return q
	}
	q.error = fmt.Errorf("sqlb: Having must be called after GroupBy")
	return q
}

func (q *UpdateQuery) Or(column string, value interface{}) *UpdateQuery {
	if q.currentTag == whereVar || q.currentTag == havingVar {
		if q.error == nil {
			if err := validateFilterColumn(column); err != nil {
				q.error = err
				return q
			}
		}
		q.clauseWrapper("or", column, value)
		return q
	}
	q.error = fmt.Errorf("sqlb: Or must be called after Where or Having")
	return q
}

func (q *UpdateQuery) GroupBy(columns ...string) *UpdateQuery {
	for _, col := range columns {
		if q.error == nil {
			if err := validateIdentifier(col); err != nil {
				q.error = err
				return q
			}
		}
	}
	cols := strings.Join(columns, ",")
	if isClauseExist(q.whereStmt, groupByVar) {
		q.whereStmt += "," + cols
	} else {
		q.whereStmt += " " + groupByVar + " " + cols
	}
	q.currentTag = groupByVar
	return q
}

func (q *UpdateQuery) Take(limit, offset int64) *UpdateQuery {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *UpdateQuery) Limit(limit int64) *UpdateQuery {
	q.limit = limit
	return q
}

func (q *UpdateQuery) Offset(offset int64) *UpdateQuery {
	q.offset = offset
	return q
}

func (q *UpdateQuery) OrderBy(columns ...string) *UpdateQuery {
	for _, col := range columns {
		if q.error == nil {
			if err := validateIdentifier(col); err != nil {
				q.error = err
				return q
			}
		}
	}
	q.orderBy = strings.Join(columns, ",")
	return q
}

func (q *UpdateQuery) Sort(sort string) *UpdateQuery {
	q.sort = sort
	return q
}

func (q *UpdateQuery) Return(fields ...string) *UpdateQuery {
	for _, f := range fields {
		if q.error == nil {
			if err := validateIdentifier(f); err != nil {
				q.error = err
				return q
			}
		}
	}
	q.returns = fields
	return q
}

func (q *UpdateQuery) Values() []any {
	return q.values
}

func (q *UpdateQuery) Stmt() string {
	return q.stmt
}

func (q *UpdateQuery) Args() []any {
	return q.values
}

func (q *UpdateQuery) Debug() string {
	return Debug(q.stmt, q.values...)
}

func (q *UpdateQuery) Error() error {
	return q.error
}

func (q *UpdateQuery) Build() *UpdateQuery {
	if q.error != nil {
		return q
	}
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

	if q.whereStmt != "" {
		q.stmt += q.whereStmt // whereStmt already starts with " "
	}

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

	if len(q.returns) > 0 {
		q.stmt += fmt.Sprintf(" returning %s;", strings.Join(q.returns, ","))
	}
	q.stmt = strings.ReplaceAll(q.stmt, beginScope, "(")
	q.stmt = strings.ReplaceAll(q.stmt, endScope, ")")
	return q
}

func (q *UpdateQuery) clause(clause string, column string, value ...interface{}) {
	queryFilter(&q.error, &q.whereStmt, &q.currentTag, &q.values, clause, column, value...)
}

func (q *UpdateQuery) clauseWrapper(clauseType string, column string, value ...interface{}) {
	filterNormalizer(&q.error, &q.whereStmt, &q.values, clauseType, column, value...)
}
