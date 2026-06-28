package sqlb

import (
	"fmt"
	"regexp"
	"strings"
)

// identifierRe validates SQL identifiers: simple names and dot-qualified names (schema.table.column).
var identifierRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)*$`)

// validateIdentifier rejects any name that contains characters outside of letters, digits,
// underscores and dots (e.g. quotes, semicolons, dashes). This prevents structural SQL injection
// in column and table names that are interpolated directly into the query string.
func validateIdentifier(s string) error {
	if s == "" {
		return fmt.Errorf("sqlb: identifier cannot be empty")
	}
	if !identifierRe.MatchString(s) {
		return fmt.Errorf("sqlb: unsafe identifier %q — only letters, digits, underscores and dots are allowed", s)
	}
	return nil
}

// validateTableExpr accepts "table", "schema.table", "table alias" and "table AS alias".
func validateTableExpr(expr string) error {
	parts := strings.Fields(expr)
	switch len(parts) {
	case 0:
		return fmt.Errorf("sqlb: empty table expression")
	case 1:
		return validateIdentifier(parts[0])
	case 2:
		if err := validateIdentifier(parts[0]); err != nil {
			return err
		}
		return validateIdentifier(parts[1])
	case 3:
		if strings.ToUpper(parts[1]) != "AS" {
			return fmt.Errorf("sqlb: invalid table expression %q", expr)
		}
		if err := validateIdentifier(parts[0]); err != nil {
			return err
		}
		return validateIdentifier(parts[2])
	default:
		return fmt.Errorf("sqlb: invalid table expression %q", expr)
	}
}

// validateFilterColumn extracts the leading identifier from a filter column string
// (which may carry internal tag prefixes and operator suffixes like "IN", "BETWEEN") and
// validates it. SQL aggregate function calls like count(*) or sum(amount) are allowed;
// only the function name (before the first `(`) is checked.
func validateFilterColumn(col string) error {
	for _, tag := range []string{betweenTag, inTag, literalTag, nullableTag} {
		col = strings.TrimPrefix(col, tag)
	}
	parts := strings.Fields(col)
	if len(parts) == 0 {
		return fmt.Errorf("sqlb: empty column name")
	}
	// For SQL functions like count(*) or sum(amount) validate only the name part.
	base := strings.SplitN(parts[0], "(", 2)[0]
	if base == "" {
		return fmt.Errorf("sqlb: empty column name")
	}
	return validateIdentifier(base)
}

// hasBetweenTag etc. are unchanged.
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

// containsOperationSymbol is kept for potential future use.
func containsOperationSymbol(str string) bool {
	regex := regexp.MustCompile(`^[<!=>]+$`)
	return regex.MatchString(str)
}
