package pgxsql

import (
	"context"
	"errors"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/host/messaging"
)

// Query - templated function for a Query
func Query[E runtime.ErrorHandler](ctx context.Context, req *Request, args ...any) (result Rows, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return nil, e.HandleWithContext(ctx, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	fn, ctx, limited = controllerApply(ctx, messaging.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if q, ok := queryExchangeCast(ctx); ok {
		var err error
		result, err = q.Query(req)
		return result, e.HandleWithContext(ctx, execLoc, err)
	}
	if dbClient == nil {
		return nil, e.HandleWithContext(ctx, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	pgxRows, err := dbClient.Query(ctx, req.BuildSql(), args...)
	if err != nil {
		return nil, e.HandleWithContext(ctx, queryLoc, recast(err))
	}
	return &proxyRows{pgxRows: pgxRows, fd: createFieldDescriptions(pgxRows.FieldDescriptions())}, runtime.NewStatusOK()
}
