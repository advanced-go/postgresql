package pgxsql

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	scanLoc = PkgPath + ":Scan"
)

// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows pgx.Rows) ([]T, *runtime.Status) {
	if rows == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, scanLoc, errors.New("invalid request: rows interface is nil"))
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
			return t, runtime.NewStatusError(runtime.StatusInvalidArgument, scanLoc, err)
		}
		values, err = rows.Values()
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusInvalidArgument, scanLoc, err)
		}
		val, err1 := s.Scan(names, values)
		if err1 != nil {
			return t, runtime.NewStatusError(runtime.StatusInvalidArgument, scanLoc, err1)
		}
		t = append(t, val)
		// Test this
		rows.Close()
	}
	return t, runtime.StatusOK()
}

func createColumnNames(fields []pgconn.FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}
