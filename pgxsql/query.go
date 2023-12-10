package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
)

const (
	queryLoc = PkgPath + ":Query"
)

// Query - function for a Query
func Query(ctx context.Context, req Request) (result pgx.Rows, status runtime.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetRequestId(ctx)
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	return queryController.Apply(ctx, req)
}

// Scrap
//var limited = false
//var fn func()
//
//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, runtime.NewStatus(runtime.StatusRateLimited)
//}
