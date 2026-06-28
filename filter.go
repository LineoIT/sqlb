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
		} else {
			*stmt += " and" // no trailing space — sub-filter adds " column"
		}
	} else {
		*stmt += " " + clause // e.g. " where", " having" — no trailing space
	}
	filterNormalizer(err, stmt, args, "", column, value...)
	*tag = clause
}

func filterNormalizer(err *error, stmt *string, args *[]interface{}, clauseType string, column string, value ...interface{}) {
	if clauseType != "" && strings.Contains(*stmt, beginScope) {
		*stmt = strings.ReplaceAll(*stmt, beginScope, "")
		clauseType = clauseType + " (" // handled by appendCond — no extra space inside paren
	}
	if hasInTag(column) {
		column = strings.ReplaceAll(column, inTag, "")
		rangeFilter(stmt, args, column, clauseType, value...)
	} else if hasNullableTag(column) {
		column = strings.ReplaceAll(column, nullableTag, "")
		nullableFilter(stmt, column, clauseType)
	} else if hasBetweenTag(column) {
		column = strings.ReplaceAll(column, betweenTag, "")
		if len(value) > 0 {
			if reflect.TypeOf(value[0]).Kind() == reflect.Slice {
				s := reflect.ValueOf(value[0])
				if s.Len() > 1 {
					betweenFilter(stmt, args, column, clauseType, s.Index(0).Interface(), s.Index(1).Interface())
				}
			} else if len(value) > 1 {
				betweenFilter(stmt, args, column, clauseType, value[0], value[1])
			}
		}
	} else if hasLiteralTag(column) {
		column = strings.ReplaceAll(column, literalTag, "")
		expressionFilter(stmt, args, column, clauseType, value[0])
	} else {
		appendCond(stmt, clauseType, fmt.Sprintf("%s=$%d", column, len(*args)+1))
		*args = append(*args, value[0])
	}
}

// appendCond writes " <clauseType> <cond>" or " <cond>" handling the special
// case where clauseType ends with "(" (scope opener) so no extra space is added
// inside the parenthesis.
func appendCond(stmt *string, clauseType, cond string) {
	if clauseType == "" {
		*stmt += " " + cond
		return
	}
	if strings.HasSuffix(clauseType, "(") {
		*stmt += " " + clauseType + cond // e.g. " or (column=$N"
	} else {
		*stmt += " " + clauseType + " " + cond // e.g. " or column=$N"
	}
}

func betweenFilter(stmt *string, args *[]interface{}, column, clauseType string, fromValue interface{}, toValue interface{}) {
	appendCond(stmt, clauseType, fmt.Sprintf("%s $%d and $%d", column, len(*args)+1, len(*args)+2))
	*args = append(*args, fromValue, toValue)
}

func nullableFilter(stmt *string, column string, clauseType string) {
	appendCond(stmt, clauseType, column)
}

func rangeFilter(stmt *string, args *[]interface{}, column, clauseType string, values ...any) {
	appendCond(stmt, clauseType, fmt.Sprintf("%s (", column))
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

func expressionFilter(stmt *string, args *[]interface{}, column, clauseType string, value interface{}) {
	if value != nil {
		appendCond(stmt, clauseType, fmt.Sprintf("%s $%d", column, len(*args)+1))
		*args = append(*args, value)
	} else {
		appendCond(stmt, clauseType, column+" null")
	}
}
