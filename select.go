package sqlb

import "strings"

type selectQuery struct {
	cols     []string
	distinct bool
}

// Select starts a SELECT query. Column expressions are trusted developer-supplied values
// (may include *, COUNT(*), aliases, etc.) and are not validated.
func Select(columns ...string) *selectQuery {
	return &selectQuery{cols: columns}
}

// Distinct adds the DISTINCT keyword.
func (s *selectQuery) Distinct() *selectQuery {
	s.distinct = true
	return s
}

func (s *selectQuery) From(tables ...string) *QueryBuilder {
	q := &QueryBuilder{args: make([]interface{}, 0)}
	for _, t := range tables {
		if err := validateTableExpr(t); err != nil {
			q.error = err
			return q
		}
	}
	cols := strings.Join(s.cols, ",")
	if cols == "" {
		cols = "*"
	}
	keyword := "select"
	if s.distinct {
		keyword = "select distinct"
	}
	q.stmt = keyword + " " + cols + " from " + strings.Join(tables, ",")
	return q
}
