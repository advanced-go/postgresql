package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
	"net/http"
)

var (
	pingLoc = PkgUri + "/Ping"
)

// Ping - templated function for pinging the database cluster
func Ping(ctx context.Context) (status *runtime.Status) {
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), PingUri, runtime.RequestId(ctx), "GET")
	defer fn()
	if limited {
		return runtime.NewStatus(runtime.StatusRateLimited)
	}
	if dbClient == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	err := dbClient.Ping(ctx)
	if err != nil {
		return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err)
	}
	return runtime.NewStatusOK()
}
