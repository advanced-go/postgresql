package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
)

const (
	queryLoc       = PkgPath + ":query"
	queryRouteName = "query"
)

// Query - function for a Query
func query(ctx context.Context, req Request) (result pgx.Rows, status runtime.Status) {
	var fn func()

	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetRequestId(ctx)
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	fn, ctx = Apply(ctx, access.NewStatusCodeClosure(&status), req.Uri(), runtime.RequestId(ctx), req.Method(), queryRouteName, queryThreshold)
	defer fn()
	rows, err := dbClient.Query(ctx, req.Sql(), req.Args())
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, queryLoc, recast(err))
	}
	return rows, runtime.StatusOK()
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
