package sqlb

import (
	"fmt"
	"strings"
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

func In[T comparable](field string, value []T) (string, []T) {
	return inTag + field + " in", value
}

func NotIn[T comparable](field string, value []T) (string, []T) {
	return inTag + field + " not in", value
}

func Between[T comparable](field string, value []T) (string, []T) {
	return betweenTag + field + " between", value
}

func IsNull(field string) (string, any) {
	return nullableTag + field + " is null", nil
}

func IsNotNull(field string) (string, any) {
	return nullableTag + field + " is not null", nil
}

func Is(field string, value any) (string, any) {
	if value == nil {
		return nullableTag + field + " is null", nil
	}
	return literalTag + field + " is", value
}

func IsNot(field string, value any) (string, any) {
	if value == nil {
		return nullableTag + field + " is not null", nil
	}
	return literalTag + field + " is not", value
}

func Ilike(field string, value any) (string, any) {
	return literalTag + field + " ilike", value
}

// StartWith : search any start at in table
func StartWith(field string, value string) (string, any) {
	return literalTag + field + " ilike", fmt.Sprintf("'%%%s'", value)
}

// EndWith : search any end at in table
func EndWith(field string, value string) (string, any) {
	return literalTag + field + " ilike", fmt.Sprintf("'%s%%'", value)
}

// Contains : search any in table
func Contains(field string, value string) (string, any) {
	return literalTag + field + " ilike", fmt.Sprintf("'%%%s%%'", value)
}

func Like[T comparable](field string, value T) (string, T) {
	return literalTag + field + " like", value
}

func Equal(field string, value any) (string, any) {
	return literalTag + field + " =", value
}

func NotEqual(field string, value any) (string, any) {
	return literalTag + field + " <>", value
}

func Greater(field string, value any) (string, any) {
	return literalTag + field + " >", value
}

func GreaterOrEqual(field string, value any) (string, any) {
	return literalTag + field + " >=", value
}

func LessOrEqual(field string, value any) (string, any) {
	return literalTag + field + " <=", value
}

func Less(field string, value any) (string, any) {
	return literalTag + field + " <", value
}

func Expression(field string, exp string, value any) (string, any) {
	var tag string
	switch strings.ToLower(exp) {
	case "is":
		return Is(field, value)
	case "is not":
		return IsNot(field, value)
	case "is null":
		return IsNull(field)
	case "is not null":
		return IsNotNull(field)
	case "between":
		tag = betweenTag
	case "in", "not in":
		tag = inTag
	default:
		tag = literalTag
	}
	return tag + field + " " + exp, value
}
