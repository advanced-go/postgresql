package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	pingLoc       = PkgPath + ":ping"
	pingRouteName = "ping"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context) (status runtime.Status) {
	if dbClient == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetRequestId(ctx)
	}
	err := dbClient.Ping(ctx)
	if err != nil {
		return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
	}
	return runtime.StatusOK()
}

// Scrap
//var limited = false
//var fn func()
//
