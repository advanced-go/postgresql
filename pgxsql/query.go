package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
)

var (
	queryLoc = PkgUri + "/Query"
)

// Query - function for a Query
func Query(ctx context.Context, req *Request) (result Rows, status *runtime.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetRequestId(ctx)
	}
	if runtime.IsDebugEnvironment() {
		if proxies, ok := runtime.IsProxyable(ctx); ok {
			if pQuery := findQueryProxy(proxies); pQuery != nil {
				var err error
				result, err = pQuery(req)
				return result, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err).SetRequestId(ctx)
			}
		}
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	pgxRows, err := dbClient.Query(ctx, BuildSql(req), req.Args)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, recast(err)).SetRequestId(ctx)
	}
	return &proxyRows{pgxRows: pgxRows, fd: createFieldDescriptions(pgxRows.FieldDescriptions())}, runtime.NewStatusOK()
}

// Scrap
//var limited = false
//var fn func()
//
//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, runtime.NewStatus(runtime.StatusRateLimited)
//}
