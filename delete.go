package sqlb

import (
	"fmt"
	"strings"
)

type DeleteQuery struct {
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

func Delete(table string) *DeleteQuery {
	return &DeleteQuery{
		stmt: "delete from " + table,
	}
}

func (q *DeleteQuery) Where(column string, value interface{}) *DeleteQuery {
	q.clause(whereVar, column, value)
	return q
}

func (q *DeleteQuery) Having(column string, value interface{}) *DeleteQuery {
	if q.currentTag == havingVar || q.currentTag == groupByVar || isClauseExist(q.whereStmt, groupByVar) {
		q.clause(havingVar, column, value)
		return q
	}
	q.error = fmt.Errorf("function `Having` should be called after group by statement")
	return q
}

func (q *DeleteQuery) Or(column string, value interface{}) *DeleteQuery {
	if q.currentTag == whereVar || q.currentTag == havingVar {
		q.clauseWrapper("or", column, value)
		return q
	}
	q.error = fmt.Errorf("function `Or` should be called after where or having statement")
	return q
}

func (q *DeleteQuery) GroupBy(columns ...string) *DeleteQuery {
	s := fmt.Sprintf(" %s %s", groupByVar, strings.Join(columns, ","))
	if isClauseExist(q.whereStmt, groupByVar) {
		q.whereStmt = strings.ReplaceAll(strings.ToLower(q.whereStmt), groupByVar, s+",")
	} else {
		q.whereStmt += s
	}
	q.currentTag = groupByVar
	return q
}

func (q *DeleteQuery) Take(limit, offset int64) *DeleteQuery {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *DeleteQuery) Limit(limit int64) *DeleteQuery {
	q.limit = limit
	return q
}

func (q *DeleteQuery) Offset(offset int64) *DeleteQuery {
	q.offset = offset
	return q
}

func (q *DeleteQuery) OrderBy(orderby ...string) *DeleteQuery {
	q.orderBy = strings.Join(orderby, ",")
	return q
}

func (q *DeleteQuery) Sort(sort string) *DeleteQuery {
	q.sort = sort
	return q
}

func (q *DeleteQuery) Return(fields ...string) *DeleteQuery {
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
		q.stmt += " " + q.whereStmt
	}

	if q.orderBy != "" {
		q.stmt += fmt.Sprintf(" order by %s", q.orderBy)
		if strings.ToLower(q.sort) == "desc" {
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
	return q
}

func (q *DeleteQuery) clause(clause string, column string, value ...interface{}) {
	queryFilter(&q.error, &q.whereStmt, &q.currentTag, &q.values, clause, column, value...)
}

func (q *DeleteQuery) clauseWrapper(clauseType string, column string, value ...interface{}) {
	filterNormalizer(&q.error, &q.whereStmt, q.currentTag, &q.values, clauseType, column, value...)
}
