package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/host"
	"github.com/go-ai-agent/core/runtime"
)

var (
	statLoc = pkgUri + "/Stat"
)

// Stat - templated function for retrieving runtime stats
func Stat[E runtime.ErrorHandler](ctx context.Context) (stat *Stats, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	fn, ctx, limited = controllerApply(ctx, host.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatus(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, e.Handle(ctx, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}
