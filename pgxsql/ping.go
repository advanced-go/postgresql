package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
	"net/http"
)

const (
	pingLoc       = PkgPath + ":ping"
	pingRouteName = "ping"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context) (status *core.Status) {
	if dbClient == nil {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	err := dbClient.Ping(ctx)
	if err != nil {
		return core.NewStatusError(http.StatusInternalServerError, err)
	}
	return core.StatusOK()
}

// Scrap
//var limited = false
//var fn func()
//
