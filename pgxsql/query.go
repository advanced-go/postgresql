package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
)

var (
	queryLoc = PkgUri + "/Query"
)

// Query - function for a Query
func Query(ctx context.Context, req Request) (result pgx.Rows, status runtime.Status) {
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
	return queryController.Apply(ctx, req)
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
