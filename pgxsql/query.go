package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
)

const (
	queryLoc = PkgPath + ":query"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status runtime.Status) {
	var fn func()
	url, ok := lookup(req.resource)

	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call : request is nil")).SetRequestId(ctx)
	}
	if !ok && dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	fn, ctx = apply(ctx, req, access.NewStatusCodeClosure(&status))
	defer fn()
	if ok {
		if len(url) == 0 {
			return nil, runtime.StatusOK()
		}
		// TO DO : create rows rom file
		return nil, runtime.StatusOK()
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
