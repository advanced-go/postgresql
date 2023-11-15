package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	statLoc = PkgUri + "/Stat"
)

// Stat - function for retrieving runtime stats
func Stat(ctx context.Context) (stat *pgxpool.Stat, status *runtime.Status) {
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil")).SetRequestId(ctx)
	}
	return dbClient.Stat(), runtime.NewStatusOK()
}

// Scrap
//var limited = false
//var fn func()

//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, runtime.NewStatus(runtime.StatusRateLimited).SetRequestId(ctx)
//}
