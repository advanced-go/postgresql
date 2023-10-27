package pgxdml

import (
	"errors"
	"strings"
)

// BuildWhere - build the []Attr based on the URL query parameters
func BuildWhere(values map[string][]string) []Attr {
	if len(values) == 0 {
		return nil
	}
	var where []Attr
	for k, v := range values {
		where = append(where, Attr{Key: k, Val: v[0]})
	}
	return where
}

// WriteWhere - build a SQL WHERE clause utilizing the given []Attr
func WriteWhere(sb *strings.Builder, terminate bool, attrs []Attr) error {
	max := len(attrs) - 1
	if max < 0 {
		return errors.New("invalid update where argument, attrs slice is empty")
	}
	sb.WriteString(Where)
	WriteWhereAttributes(sb, attrs)
	if terminate {
		sb.WriteString(";")
	}
	return nil
}

// WriteWhereAttributes - build a SQL statement only containing the []Attr conditionals
func WriteWhereAttributes(sb *strings.Builder, attrs []Attr) error {
	max := len(attrs) - 1
	if max < 0 {
		return errors.New("invalid update where argument, attrs slice is empty")
	}
	for i, attr := range attrs {
		s, err := FmtAttr(attr)
		if err != nil {
			return err
		}
		sb.WriteString(s)
		if i < max {
			sb.WriteString(And)
		}
	}
	//sb.WriteString(";")
	return nil
}
