package pgxsql

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
	Rows([]T) [][]any
}

// Rows - templated function for creating rows
func Rows[T Scanner[T]](entries []T) ([][]any, *core.Status) {
	if len(entries) == 0 {
		return nil, core.StatusNotFound()
	}
	var t T
	return t.Rows(entries), core.StatusOK()
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows pgx.Rows) ([]T, *core.Status) {
	if rows == nil || rows.CommandTag().RowsAffected() == 0 {
		return nil, core.StatusNotFound() //core.NewStatusError(core.StatusInvalidArgument, errors.New("invalid request: rows interface is nil"))
	}
	var s T
	var t []T
	var err error
	var values []any

	defer rows.Close()
	names := createColumnNames(rows.FieldDescriptions())
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return t, core.NewStatusError(core.StatusInvalidArgument, err)
		}
		values, err = rows.Values()
		if err != nil {
			return t, core.NewStatusError(core.StatusInvalidArgument, err)
		}
		val, err1 := s.Scan(names, values)
		if err1 != nil {
			return t, core.NewStatusError(core.StatusInvalidArgument, err1)
		}
		t = append(t, val)
		// Test this
		//rows.Close()
	}
	if len(t) == 0 {
		return t, core.StatusNotFound()
	}
	return t, core.StatusOK()
}

func createColumnNames(fields []pgconn.FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}
