package pgxdml

import (
	"errors"
	"strings"
)

const (
	whereClause = "{where}"
)

// ExpandSelect - given a template, expand the template to build a WHERE clause if configured
func ExpandSelect(template string, where []Attr) (string, error) {
	if template == "" {
		return template, errors.New("template is empty")
	}
	pos := strings.Index(template, whereClause)
	if pos == -1 {
		return template, nil
	}
	sb := strings.Builder{}
	if len(where) == 0 {
		sb.WriteString(template[:pos])
		sb.WriteString(template[pos+len(whereClause)+1:])
		return sb.String(), nil
	}
	sb.WriteString(template[:pos])
	sb.WriteString("\nWHERE ")
	WriteWhereAttributes(&sb, where)
	sb.WriteString(template[pos+len(whereClause):])
	return sb.String(), nil
}
