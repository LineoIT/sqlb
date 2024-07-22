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

func (q *QueryBuilder) OrderBy(orderby ...string) *QueryBuilder {
	q.orderBy = strings.Join(orderby, ",")
	return q
}

func (q *QueryBuilder) Sort(sort string) *QueryBuilder {
	q.sort = sort
	return q
}

func (q *QueryBuilder) Or(column string, value interface{}) *QueryBuilder {
	if q.currentTag == whereVar || q.currentTag == havingVar {
		q.normalizeFilter("or", column, value)
		return q
	}
	q.error = fmt.Errorf("function `Or` should be called after where or having statement")
	return q
}

func (q *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	s := fmt.Sprintf(" %s %s", groupByVar, strings.Join(columns, ","))
	if isClauseExist(q.stmt, groupByVar) {
		q.stmt = strings.ReplaceAll(strings.ToLower(q.stmt), groupByVar, s+",")
	} else {
		q.stmt += s
	}
	q.currentTag = groupByVar
	return q
}

func (q *QueryBuilder) Having(column string, value interface{}) *QueryBuilder {
	if q.currentTag == havingVar || q.currentTag == groupByVar || isClauseExist(q.stmt, groupByVar) {
		q.filter(havingVar, column, value)
		return q
	}
	q.error = fmt.Errorf("function `Having` should be called after group by statement")
	return q
}

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

func (q *QueryBuilder) WhereRaw(raw string) *QueryBuilder {
	if strings.Contains(strings.ToLower(q.stmt), "where") {
		q.stmt += " and "
	} else {
		q.stmt += " where "
		q.currentTag = whereVar
	}
	q.stmt += " " + raw
	return q
}

func (q *QueryBuilder) OrRaw(raw string) *QueryBuilder {
	if strings.Contains(strings.ToLower(q.stmt), "where") {
		q.stmt += " or "
	} else {
		q.stmt += " where "
		q.currentTag = whereVar
	}
	q.stmt += " " + raw
	return q
}

func (q *QueryBuilder) Contains(col string, value string) *QueryBuilder {
	c, v := Ilike(col, "%"+value+"%")
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) StartWith(col string, value string) *QueryBuilder {
	// q.Where(Ilike(col, value+"%"))
	c, v := Ilike(col, value+"%")
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) EndWith(col string, value string) *QueryBuilder {
	// q.Where(Ilike(col, "%"+value))
	c, v := Ilike(col, "%"+value)
	q.matchAny(c, v)
	return q
}

func (q *QueryBuilder) matchAny(col string, value any) {
	// orIndex := strings.LastIndex(strings.ToLower(q.stmt), " or ")
	// andIndex := strings.LastIndex(strings.ToLower(q.stmt), " and ")
	// if orIndex > andIndex {
	// 	q.Or(col, value)
	// } else
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
	filterNormalizer(&q.error, &q.stmt, q.currentTag, &q.args, clauseType, column, value...)
}
