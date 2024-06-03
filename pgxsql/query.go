package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
	"time"
)

const (
	queryLoc = PkgPath + ":query"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status *core.Status) {
	if req == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call : request is nil"))
	}
	if dbClient == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call: dbClient is nil"))
	}
	var err error
	ctx1, cancel := setTimeout(ctx, req)
	if cancel != nil {
		defer cancel()
	}
	var start = time.Now().UTC()
	if req.queryFunc != nil {
		rows, err = req.queryFunc(ctx1, buildSql(req), req.args)
	} else {
		rows, err = dbClient.Query(ctx1, buildSql(req), req.args)
	}
	if err != nil {
		status = core.NewStatusError(core.StatusIOError, recast(err))
	} else {
		status = core.StatusOK()
	}
	// TODO : determine if there was a timeout
	log(start, time.Since(start), req, status, "")
	return rows, status
}

func queryFunc(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {

	return nil, nil
}

// Scrap
//url, override := lookup.Value(req.resource)
//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
//if override {
//	// TO DO : create rows from file
//	return io2.New[pgx.Rows](url, nil)
//}
//var limited = false
//var fn func()
//
//fn, ctx, limited = controllerApply(ctx, startup.NewStatusCode(&status), req.Uri, core.ContextRequestId(ctx), "GET")
//defer fn()
//if limited {
//	return nil, core.NewStatus(core.StatusRateLimited)
//}
