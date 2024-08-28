package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
	"time"
)

// Query - function for a Query
func query(ctx context.Context, req *request) (rows pgx.Rows, status *core.Status) {
	if req == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call : request is nil"))
	}
	if dbClient == nil && req.queryFunc == nil {
		return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL database query call: dbClient is nil"))
	}
	var err error
	ctx = req.setTimeout(ctx)

	var start = time.Now().UTC()
	if req.queryFunc != nil {
		rows, err = req.queryFunc(ctx, buildSql(req), req)
	} else {
		rows, err = dbClient.Query(ctx, buildSql(req), req.args)
	}
	if err != nil {
		status = core.NewStatusError(core.StatusIOError, recast(err))
	} else {
		status = core.StatusOK()
	}
	// TODO : determine if there was a timeout
	reasonCode := ""
	access.Log(access.EgressTraffic, start, time.Since(start), req, status, access.Routing{From: req.From(), Route: req.routeName, To: ""}, access.Controller{Timeout: req.duration, RateLimit: -1, RateBurst: -1, Code: reasonCode})
	return rows, status
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
