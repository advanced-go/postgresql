package pgxdml

import (
	"strings"
)

// WriteDelete - build a SQL delete statement with a WHERE clause
func WriteDelete(sql string, where []Attr) (string, error) {
	var sb strings.Builder

	sb.WriteString(sql)
	sb.WriteString("\n")
	err := WriteWhere(&sb, true, where)
	return sb.String(), err
}
