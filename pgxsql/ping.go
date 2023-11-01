package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

var (
	pingLoc = PkgUri + "/Ping"
)

// Ping - function for pinging the database cluster
func Ping(ctx context.Context) (status *runtime.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetRequestId(ctx)
	}
	err := dbClient.Ping(ctx)
	if err != nil {
		return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
	}
	return runtime.NewStatusOK()
}

// Scrap
//var limited = false
//var fn func()
//
