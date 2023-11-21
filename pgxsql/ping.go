package pgxsql

import (
	"context"
	"github.com/advanced-go/core/runtime"
)

const (
	pingLoc = PkgUri + "/Ping"
)

// Ping - function for pinging the database cluster
func Ping(ctx context.Context) (status runtime.Status) {
	////if dbClient == nil {
	//	return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetRequestId(ctx)
	//}
	//err := dbClient.Ping(ctx)
	//if err != nil {
	//	return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
	//}
	//return runtime.NewStatusOK()
	return pingController.Apply(ctx)
}

// Scrap
//var limited = false
//var fn func()
//
