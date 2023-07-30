package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
)

// Query - templated function for a Query
func Query[E runtime.ErrorHandler](ctx context.Context, req *Request) (result Rows, status *runtime.Status) {
	var e E
	//var h H
	var limited = false
	//var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return nil, e.Handle(ctx, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	//	fn, ctx, limited = h.Apply(ctx, host.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
	//	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	//if q, ok := queryExchangeCast(ctx); ok {
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		if pQuery := findQueryProxy(proxies); pQuery != nil {
			var err error
			result, err = pQuery(req)
			return result, e.Handle(ctx, execLoc, err)
		}
	}
	if dbClient == nil {
		return nil, e.Handle(ctx, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	pgxRows, err := dbClient.Query(ctx, BuildSql(req), req.Args)
	if err != nil {
		return nil, e.Handle(ctx, queryLoc, recast(err))
	}
	return &proxyRows{pgxRows: pgxRows, fd: createFieldDescriptions(pgxRows.FieldDescriptions())}, runtime.NewStatusOK()
}
