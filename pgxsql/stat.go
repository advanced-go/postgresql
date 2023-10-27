package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
)

var (
	statLoc = PkgUri + "/Stat"
)

// Stat - function for retrieving runtime stats
func Stat(ctx context.Context) (stat *Stats, status *runtime.Status) {
	var limited = false
	var fn func()

	fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return nil, runtime.NewStatus(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument).SetRequestId(ctx)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}
