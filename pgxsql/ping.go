package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/host"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/resiliency/controller"
)

var (
	pingLoc = pkgPath + "/stat"
)

// Ping - templated function for pinging the database cluster
func Ping[E runtime.ErrorHandler, H controller.Handler](ctx context.Context) (status *runtime.Status) {
	var e E
	var h H
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = h.Apply(ctx, host.NewStatusCode(&status), PingUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return e.Handle(ctx, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return e.Handle(ctx, pingLoc, dbClient.Ping(ctx))
}
