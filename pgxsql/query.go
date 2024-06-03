package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
)

const (
	queryLoc = PkgPath + ":query"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status *core.Status) {
	//url, override := lookup.Value(req.resource)
	var newCtx context.Context

	if req == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call : request is nil"))
	}
	//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
	//if override {
	//	// TO DO : create rows from file
	//	return io2.New[pgx.Rows](url, nil)
	//}
	if dbClient == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call: dbClient is nil"))
	}
	var err error
	rows, err = dbClient.Query(newCtx, buildSql(req), req.args)
	if err != nil {
		return nil, core.NewStatusError(core.StatusIOError, recast(err))
	}
	return rows, core.StatusOK()
}

// Scrap
//var limited = false
//var fn func()
//
//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), req.Uri, core.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, core.NewStatus(core.StatusRateLimited)
//}
