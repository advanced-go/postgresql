package pgxdml

import (
	"errors"
	"strings"
)

/*
INSERT INTO table_name (column_list) VALUES
    (value_list_1),
    (value_list_2),
    ...
    (value_list_n);

*/

// NewInsertValues - create a 2-dimensional array of any out of a slice of any
func NewInsertValues(v []any) [][]any {
	if len(v) == 0 {
		return nil
	}
	var values [][]any
	return append(values, v)
}

// WriteInsert - build a SQL insert statement with a VALUES list
func WriteInsert(sql string, values [][]any) (string, error) {
	sb := strings.Builder{}

	sb.WriteString(sql)
	sb.WriteString("\n")
	for i, val := range values {
		if i > 0 {
			sb.WriteString(",\n")
		}
		err := WriteInsertValues(&sb, val)
		if err != nil {
			return sb.String(), err
		}
	}
	sb.WriteString(";\n")
	return sb.String(), nil
}

// WriteInsertValues - build the values list of a SQL insert statement
func WriteInsertValues(sb *strings.Builder, values []any) error {
	max := len(values) - 1
	if max < 0 {
		return errors.New("invalid insert argument, values slice is empty")
	}
	sb.WriteString("(")
	for i, v := range values {
		s, err := FmtValue(v)
		if err != nil {
			return err
		}
		sb.WriteString(s)
		if i < max {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")
	return nil
}
