package sqlb

import "fmt"

func (q *QueryBuilder) Join(table string) *QueryBuilder {
	q.currentTag = joinVar
	q.stmt += " join " + table
	return q
}

func (q *QueryBuilder) InnerJoin(table string) *QueryBuilder {
	q.currentTag = joinVar
	q.stmt += " inner join " + table
	return q
}

func (q *QueryBuilder) LeftJoin(table string) *QueryBuilder {
	q.currentTag = joinVar
	q.stmt += " left join " + table
	return q
}

func (q *QueryBuilder) RightJoin(table string) *QueryBuilder {
	q.currentTag = joinVar
	q.stmt += " right join " + table
	return q
}

func (q *QueryBuilder) On(foreignKey, referenceKey string) *QueryBuilder {
	q.currentTag = joinOnVar
	q.stmt += fmt.Sprintf(" on %s=%s", foreignKey, referenceKey)
	return q
}
