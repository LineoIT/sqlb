package sqlb

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func recFuncValue(fv ValueFunc, argIndex int) (string, any) {
	var castype string
	if fv.cast != "" {
		castype = "::" + fv.cast
	}
	val, ok := fv.value.(ValueFunc)
	s := fmt.Sprintf("%s($%v%s,%v)", fv.fun, argIndex, castype, fv.alternative)
	value := fv.value
	if ok {
		t, v := recFuncValue(val, argIndex)
		value = v
		s = strings.ReplaceAll(s, fmt.Sprintf("$%v", argIndex), t)
	}
	return s, value
}

func CleanSQL(query string) string {
	return strings.ReplaceAll(strings.Trim(strings.ReplaceAll(
		strings.ReplaceAll(query, "\t", ""), "\n", " "), " "), "  ", " ")
}

func Debug(query string, args ...interface{}) string {
	s := CleanSQL(query)
	for k, v := range args {
		s = strings.Replace(s, fmt.Sprintf("$%d", k+1), fmt.Sprint(v), -1)
	}
	return s
}

func mergeSliceValue(arr []any) []any {
	var list []any
	for _, t := range arr {
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				list = append(list, s.Index(i).Interface())
			}
		default:
			list = append(list, t)
		}
	}
	return list
}

func isInclusingClause(column string) bool {
	return strings.Count(column, inVar) == 1 || strings.Count(column, notInVar) == 1
}

func isNullableClause(column string) bool {
	return strings.Count(column, isNullVar) == 1 || strings.Count(column, isNotNullVar) == 1
}

func isIntervalClause(column string) bool {
	return strings.Count(column, betweenVar) == 1
}

func isExpressionClause(column string) bool {
	return strings.Count(column, isVar) == 1 ||
		strings.Count(column, isNotVar) == 1 ||
		strings.Count(column, likeVar) == 1 ||
		strings.Count(column, ilikeVar) == 1 ||
		isValidColumnExpression(column)
}

// Pattern to match alphanumeric characters and underscore
func isAllowedColumnName(str string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*[a-zA-Z0-9]$`)
	return regex.MatchString(str)
}

// Pattern to match allowed characters
func isAllowedExpression(str string) bool {
	regex := regexp.MustCompile(`^[<!=>]+$`)
	return regex.MatchString(str)
}

func isValidColumnExpression(str string) bool {
	words := strings.Fields(str)
	if len(words) == 2 {
		return isAllowedColumnName(words[0]) && isAllowedExpression(words[1])
	}
	return false
}
