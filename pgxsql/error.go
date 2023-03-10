package pgxsql

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

/*
code
message
detail
hint
position
internalPosition
internalQuery
where
SchemaName
TableName
ColumnName
DataTypeName
ConstraintName

*/
func recast(err error) error {
	if err == nil {
		return err
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		msg := fmt.Sprintf("serverity:%v, code:%v, message:%v, position:%v, SQLState:%v",
			pgErr.Severity, pgErr.Code, pgErr.Message, pgErr.Position, pgErr.SQLState())
		//, pgErr.Line)
		//pgErr.Code,pgErr.ColumnName,pgErr.ConstraintName,pgErr.DataTypeName,pgErr.Message,
		//pgErr.Detail,pgErr.File,pgErr.Hint,pgErr.InternalPosition,pgErr.Line
		return errors.New(msg)
	}
	return err
}
