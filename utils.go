package sqlb

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var argPlaceholderRe = regexp.MustCompile(`\$(\d+)`)

// renumberArgs shifts all $N placeholders in stmt by offset.
// Used when composing sub-queries or CTEs so arg indices don't collide.
func renumberArgs(stmt string, offset int) string {
	if offset == 0 {
		return stmt
	}
	return argPlaceholderRe.ReplaceAllStringFunc(stmt, func(match string) string {
		n, _ := strconv.Atoi(match[1:])
		return fmt.Sprintf("$%d", n+offset)
	})
}

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
	// Substitute from highest index first so $1 never clobbers $10, $11, etc.
	for k := len(args) - 1; k >= 0; k-- {
		s = strings.Replace(s, fmt.Sprintf("$%d", k+1), fmt.Sprint(args[k]), -1)
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
