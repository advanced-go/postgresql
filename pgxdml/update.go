package pgxdml

import (
	"errors"
	"strings"
)

const (
	Where = "WHERE "
	And   = " AND "
	Set   = "SET "
)

/*
UPDATE table_name
SET column1 = value1,
    column2 = value2,
    ...
WHERE condition;

*/

// WriteUpdate - build a SQL update statement, including SET and WHERE clauses
func WriteUpdate(sql string, attrs []Attr, where []Attr) (string, error) {
	var sb strings.Builder

	sb.WriteString(sql)
	sb.WriteString("\n")
	err := WriteUpdateSet(&sb, attrs)
	if err != nil {
		return "", err
	}
	err = WriteWhere(&sb, true, where)
	return sb.String(), err
}

// WriteUpdateSet - build a SQL set clause
func WriteUpdateSet(sb *strings.Builder, attrs []Attr) error {
	max := len(attrs) - 1
	if max < 0 {
		return errors.New("invalid update set argument, attrs slice is empty")
	}
	sb.WriteString(Set)
	for i, attr := range attrs {
		s, err := FmtAttr(attr)
		if err != nil {
			return err
		}
		sb.WriteString(s)
		if i < max {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n")
	return nil
}
