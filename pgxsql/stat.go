package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/controller"
	"github.com/go-ai-agent/core/resource"
	"github.com/go-ai-agent/core/runtime"
)

var (
	statLoc = pkgPath + "/stat"
)

// Stat - templated function for retrieving runtime stats
func Stat[E runtime.ErrorHandler, H controller.Handler](ctx context.Context) (stat *Stats, status *runtime.Status) {
	var e E
	var h H
	var limited = false
	var fn func()

	fn, ctx, limited = h.Apply(ctx, resource.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.Handle(ctx, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}
