package sqlb

import (
	"regexp"
	"strings"
)

// Pattern to match allowed characters
func containsOperationSymbol(str string) bool {
	regex := regexp.MustCompile(`^[<!=>]+$`)
	return regex.MatchString(str)
}

func hasBetweenTag(col string) bool {
	return strings.HasPrefix(col, betweenTag)
}

func hasInTag(col string) bool {
	return strings.HasPrefix(col, inTag)
}

func hasLiteralTag(col string) bool {
	return strings.HasPrefix(col, literalTag)
}

func hasNullableTag(col string) bool {
	return strings.HasPrefix(col, nullableTag)
}

func isClauseExist(stmt, clause string) bool {
	stmt = strings.ToLower(stmt)
	return strings.LastIndex(stmt, clause) > strings.LastIndex(stmt, "from")
}
