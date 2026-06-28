package sqlb

import "fmt"

func (q *QueryBuilder) Join(table string) *QueryBuilder {
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.currentTag = joinVar
	q.stmt += " join " + table
	return q
}

func (q *QueryBuilder) InnerJoin(table string) *QueryBuilder {
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.currentTag = joinVar
	q.stmt += " inner join " + table
	return q
}

func (q *QueryBuilder) LeftJoin(table string) *QueryBuilder {
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.currentTag = joinVar
	q.stmt += " left join " + table
	return q
}

func (q *QueryBuilder) RightJoin(table string) *QueryBuilder {
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.currentTag = joinVar
	q.stmt += " right join " + table
	return q
}

func (q *QueryBuilder) FullJoin(table string) *QueryBuilder {
	if err := validateTableExpr(table); err != nil {
		q.error = err
		return q
	}
	q.currentTag = joinVar
	q.stmt += " full join " + table
	return q
}

// On adds an ON clause for a join. Both keys must be valid identifiers (table.column).
func (q *QueryBuilder) On(foreignKey, referenceKey string) *QueryBuilder {
	if q.error == nil {
		if err := validateIdentifier(foreignKey); err != nil {
			q.error = err
			return q
		}
		if err := validateIdentifier(referenceKey); err != nil {
			q.error = err
			return q
		}
	}
	q.currentTag = joinOnVar
	q.stmt += fmt.Sprintf(" on %s=%s", foreignKey, referenceKey)
	return q
}
