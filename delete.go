package sqlb

import (
	"fmt"
	"strings"
)

type DeleteQuery struct {
	values        []any
	returns       []string
	stmt          string
	whereStmt     string
	currentTag    string
	error         error
	limit, offset int64
	orderBy, sort string
}

func Delete(table string) *DeleteQuery {
	q := &DeleteQuery{}
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.stmt = "delete from " + table
	return q
}

func (q *DeleteQuery) Where(column string, value interface{}) *DeleteQuery {
	if q.error == nil {
		if err := validateFilterColumn(column); err != nil {
			q.error = err
			return q
		}
	}
	q.clause(whereVar, column, value)
	return q
}

func (q *DeleteQuery) Having(column string, value interface{}) *DeleteQuery {
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

func (q *DeleteQuery) Or(column string, value interface{}) *DeleteQuery {
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

func (q *DeleteQuery) Return(fields ...string) *DeleteQuery {
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

func (q *DeleteQuery) Values() []any {
	return q.values
}

func (q *DeleteQuery) Stmt() string {
	return q.stmt
}

func (q *DeleteQuery) Debug() string {
	return Debug(q.stmt, q.values...)
}

func (q *DeleteQuery) Error() error {
	return q.error
}

func (q *DeleteQuery) Build() *DeleteQuery {
	if q.error != nil {
		return q
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

func (q *DeleteQuery) clause(clause string, column string, value ...interface{}) {
	queryFilter(&q.error, &q.whereStmt, &q.currentTag, &q.values, clause, column, value...)
}

func (q *DeleteQuery) clauseWrapper(clauseType string, column string, value ...interface{}) {
	filterNormalizer(&q.error, &q.whereStmt, &q.values, clauseType, column, value...)
}
