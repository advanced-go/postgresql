package pgxsql

import (
	"errors"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5/pgxpool"
)

func stat() (*pgxpool.Stat, *core.Status) {
	if dbClient == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL stat call : dbClient is nil"))
	}
	return dbClient.Stat(), core.StatusOK()
}

// Scrap
//var limited = false
//var fn func()

//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), StatUri, core.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, core.NewStatus(core.StatusRateLimited).SetRequestId(ctx)
//}
