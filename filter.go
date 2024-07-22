package sqlb

import (
	"fmt"
	"reflect"
	"strings"
)

func queryFilter(err *error, stmt *string, tag *string, args *[]interface{}, clause string, column string, value ...interface{}) {

	if isClauseExist(*stmt, clause) {
		if strings.Contains(*stmt, beginScope) {
			*stmt = strings.ReplaceAll(*stmt, beginScope, "")
			*stmt += " and ("
		}
	} else {
		*stmt += " " + clause
	}
	filterNormalizer(err, stmt, *tag, args, "", column, value...)
	*tag = clause
}

func filterNormalizer(err *error, stmt *string, tag string, args *[]interface{}, clauseType string, column string, value ...interface{}) {
	if clauseType != "" && strings.Contains(*stmt, beginScope) {
		*stmt = strings.ReplaceAll(*stmt, beginScope, "")
		clauseType = " " + clauseType + " ("
	}
	if hasInTag(column) {
		column = strings.ReplaceAll(column, inTag, "")
		rangeFilter(stmt, tag, args, column, clauseType, value...)
	} else if hasNullableTag(column) {
		column = strings.ReplaceAll(column, nullableTag, "")
		nullableFilter(stmt, tag, column, clauseType)
	} else if hasBetweenTag(column) {
		column = strings.ReplaceAll(column, betweenTag, "")
		if len(value) > 0 {
			if reflect.TypeOf(value[0]).Kind() == reflect.Slice {
				s := reflect.ValueOf(value[0])
				if s.Len() > 1 {
					fromValue := s.Index(0).Interface()
					toValue := s.Index(1).Interface()
					betweenFilter(stmt, tag, args, column, clauseType, fromValue, toValue)
				}
			} else if len(value) > 1 {
				betweenFilter(stmt, tag, args, column, clauseType, value[0], value[1])
			}
		}
	} else if hasLiteralTag(column) {
		column = strings.ReplaceAll(column, literalTag, "")
		expressionFilter(stmt, tag, args, column, clauseType, value[0])
	} else {
		*stmt += fmt.Sprintf(" %s %s=$%d", clauseType, column, len(*args)+1)
		*args = append(*args, value[0])
	}
}

func betweenFilter(stmt *string, tag string, args *[]interface{}, column, clauseType string, fromValue interface{}, toValue interface{}) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += " " + clauseType
	} else {
		*stmt += " " + tag
	}
	*stmt += fmt.Sprintf(" %s $%d and $%d", column, len(*args)+1, len(*args)+2)
	*args = append(*args, fromValue, toValue)
}

func nullableFilter(stmt *string, tag, column string, clauseType string) {
	if strings.Count(*stmt, tag) > 0 {
		*stmt += fmt.Sprintf(" %s %s", clauseType, column)
	} else {
		*stmt += fmt.Sprintf(" %s %s", tag, column)
	}
}

func rangeFilter(stmt *string, tag string, args *[]interface{}, column, clauseType string, values ...any) {
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

func expressionFilter(stmt *string, tag string, args *[]interface{}, column, clauseType string, value interface{}) {
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
