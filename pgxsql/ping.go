package pgxsql

import (
	"context"
	"errors"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/host/messaging"
)

var (
	pingLoc = pkgPath + "/stat"
)

// Ping - templated function for pinging the database cluster
func Ping[E runtime.ErrorHandler](ctx context.Context) (status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = controllerApply(ctx, messaging.NewStatusCode(&status), PingUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return e.HandleWithContext(ctx, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return e.HandleWithContext(ctx, pingLoc, dbClient.Ping(ctx))
}
