package sqlb

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

func Eq[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " =", value
}

func NotEq[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " <>", value
}

func In[T comparable](field string, value []T) (string, []T) {
	checkValidColumn(field)
	return field + " in", value
}

func NotIn[T comparable](field string, value []T) (string, []T) {
	checkValidColumn(field)
	return field + " not in", value
}

func Between[T comparable](field string, value []T) (string, []T) {
	checkValidColumn(field)
	return field + " between", value
}

func IsNull(field string) string {
	checkValidColumn(field)
	return field + " is null"
}

func IsNotNull(field string) string {
	checkValidColumn(field)
	return field + " is not null"
}

func Is[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " is", value
}

func IsNot[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " is not", value
}

func Ilike[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " ilike", value
}

func Like[T comparable](field string, value T) (string, T) {
	checkValidColumn(field)
	return field + " like", value
}

func Expression[T comparable](field string, exp string, value T) (string, T) {
	checkValidColumn(field)
	if !isAllowedExpression(exp) {
		panic(exp + " is not allowed character")
	}
	return field + " " + exp, value
}

func checkValidColumn(field string) {
	if !isAllowedColumnName(field) {
		panic(field + " is not allowed character")
	}
}
