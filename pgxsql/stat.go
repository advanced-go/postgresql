package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	statLoc = PkgPath + ":stat"
)

func stat(ctx context.Context) (*pgxpool.Stat, *runtime.Status) {
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, statLoc, errors.New("error on PostgreSQL stat call : dbClient is nil"))
	}
	return dbClient.Stat(), runtime.StatusOK()
}

// Scrap
//var limited = false
//var fn func()

//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, runtime.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, runtime.NewStatus(runtime.StatusRateLimited).SetRequestId(ctx)
//}
