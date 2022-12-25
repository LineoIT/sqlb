package sqlb

import (
	"fmt"
	"strings"
	"reflect"
)

type ValueFunc struct {
	value       any
	alternative any
	fun         string
	cast        string
}

func Coalesce(value any, alternative any, cast ...string) ValueFunc {
	v := ValueFunc{
		value:       value,
		fun:         "coalesce",
		alternative: alternative,
	}
	if len(cast) > 0 {
		v.cast = cast[0]
	}
	return v
}

func Nullif(value any, alternative any, cast ...string) ValueFunc {
	v := ValueFunc{
		value:       value,
		fun:         "nullif",
		alternative: alternative,
	}
	if len(cast) > 0 {
		v.cast = cast[0]
	}
	return v
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