package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
)

const (
	queryLoc = PkgPath + ":query"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status runtime.Status) {
	var fn func()
	urls := lookup(req.resource)

	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetRequestId(ctx)
	}
	if urls == nil && dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	fn, ctx = apply(ctx, req, &status)
	defer fn()
	if urls != nil {
		// TO DO : create rows from file
		return io2.ReadResults[pgx.Rows](urls)
	}
	var err error
	rows, err = dbClient.Query(ctx, buildSql(req), req.args)
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
