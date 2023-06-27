package sqlb

import (
	"fmt"
	"reflect"
	"strings"
)

func commonClause(err *error, stmt *string, tag *string, args *[]interface{}, clause string, column string, value ...interface{}) {
	if strings.Count(*stmt, clause) > 0 {
		*stmt += " and"
	} else {
		*stmt += " " + clause
	}
	clauseWrapper(err, stmt, *tag, args, "", column, value...)
	//q.clauseWrapper("", column, value...)
	//q.currentTag = clause
	*tag = clause
}

func clauseWrapper(err *error, stmt *string, tag string, args *[]interface{}, clauseType string, column string, value ...interface{}) {
	if isInclusingClause(column) {
		rangeClause(stmt, tag, args, column, clauseType, value...)
		// q.rangeClause(column, clauseType, value...)
	} else if isNullableClause(column) {
		nullableClause(stmt, tag, column, clauseType)
		//	q.nullableClause(column, clauseType)
	} else if isIntervalClause(column) {
		if len(value) > 0 {
			if reflect.TypeOf(value[0]).Kind() == reflect.Slice {
				s := reflect.ValueOf(value[0])
				if s.Len() > 1 {
					fromValue := s.Index(0).Interface()
					toValue := s.Index(1).Interface()
					between(stmt, tag, args, column, clauseType, fromValue, toValue)
					// q.between(column, clauseType, fromValue, toValue)
				}
			} else if len(value) > 1 {
				between(stmt, tag, args, column, clauseType, value[0], value[1])
				// q.between(column, clauseType, value[0], value[1])
			}
		}
	} else if isExpressionClause(column) {
		expressionClause(stmt, tag, args, column, clauseType, value[0])
		//q.expressionClause(column, clauseType, value[0])
	} else if isAllowedColumnName(column) {
		// q._Stmt += fmt.Sprintf(" %s %s=$%d", clauseType, column, len(q.args)+1)
		*stmt += fmt.Sprintf(" %s %s=$%d", clauseType, column, len(*args)+1)
		//q.args = append(q.args, value[0])
		*args = append(*args, value[0])
	} else {
		*err = fmt.Errorf("expression is not allowed")
	}
}

func between(stmt *string, tag string, args *[]interface{}, column, clauseType string, fromValue interface{}, toValue interface{}) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += " " + clauseType
	} else {
		*stmt += " " + tag
	}
	*stmt += fmt.Sprintf(" %s $%d and $%d", column, len(*args)+1, len(*args)+2)
	*args = append(*args, fromValue, toValue)
}

func nullableClause(stmt *string, tag, column string, clauseType string) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += fmt.Sprintf(" %s %s", clauseType, column)
	} else {
		*stmt += fmt.Sprintf(" %s %s", tag, column)
	}
}

func rangeClause(stmt *string, tag string, args *[]interface{}, column, clauseType string, values ...any) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += fmt.Sprintf(" %s %s (", clauseType, column)
	} else {
		*stmt += fmt.Sprintf(" %s %s (", tag, column)
	}
	values = mergeSliceValue(values)
	for k, v := range values {
		*stmt += fmt.Sprintf("$%d", len(*args)+1)
		if k < len(values)-1 {
			*stmt += ","
		}
		*args = append(*args, v)
	}
	*stmt += ")"
}

func expressionClause(stmt *string, tag string, args *[]interface{}, column, clauseType string, value interface{}) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += " " + clauseType
	} else {
		*stmt += " " + tag
	}
	if value != nil {
		*stmt += fmt.Sprintf(" %s $%d", column, len(*args)+1)
		*args = append(*args, value)
	} else {
		*stmt += fmt.Sprintf(" %s null", column)
	}
}
