package sqlb

import (
	"regexp"
	"strings"
)

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
	regex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.:]*[a-zA-Z0-9]$`)
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

func checkValidColumn(field string) {
	if !isAllowedColumnName(field) {
		panic(field + " is not allowed character")
	}
}
