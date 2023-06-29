package sqlb

import "strings"

type selectQuery struct {
	cols []string
}

func Select(columns ...string) *selectQuery {
	return &selectQuery{
		cols: columns,
	}
}

func (s *selectQuery) From(tables ...string) *QueryBuilder {
	cols := strings.Join(s.cols, ",")
	if cols == "" {
		cols = "*"
	}
	baseQuery := "select " + cols + " from " + strings.Join(tables, ",")
	return &QueryBuilder{
		stmt: baseQuery,
		args: make([]interface{}, 0),
	}
}
