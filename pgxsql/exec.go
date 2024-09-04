package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"time"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status *core.Status) {
	reasonCode := ""
	if req == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : request is nil"))
	}
	if dbClient == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : dbClient is nil"))
	}
	ctx = req.setTimeout(ctx)

	var start = time.Now().UTC()

	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		status = core.NewStatusError(core.StatusTxnBeginError, err0)
		// TODO : determine if there was a timeout
		access.Log(access.EgressTraffic, start, time.Since(start), req, status, access.Routing{From: req.From(), Route: req.routeName, To: ""}, access.Controller{Timeout: req.duration, RateLimit: 0, RateBurst: 0, Code: reasonCode})
		return tag, status
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err := dbClient.Exec(ctx, buildSql(req), req.args)
	if err != nil {
		status = core.NewStatusError(core.StatusInvalidArgument, recast(err))
		// TODO : determine if there was a timeout
		access.Log(access.EgressTraffic, start, time.Since(start), req, status, access.Routing{From: req.From(), Route: req.routeName, To: ""}, access.Controller{Timeout: req.duration, RateLimit: 0, RateBurst: 0, Code: reasonCode})
		return newCmdTag(cmd), status
	}
	err = txn.Commit(ctx)
	if err != nil {
		status = core.NewStatusError(core.StatusTxnCommitError, err)
	} else {
		status = core.StatusOK()
	}
	// TODO : determine if there was a timeout
	access.Log(access.EgressTraffic, start, time.Since(start), req, status, access.Routing{From: req.From(), Route: req.routeName, To: ""}, access.Controller{Timeout: req.duration, RateLimit: 0, RateBurst: 0, Code: reasonCode})
	return newCmdTag(cmd), core.StatusOK()
}

// scrap
//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
//if override {
//	return io2.New[CommandTag](url, nil)
//}
//
//url, override := lookup.Value(req.resource)
//var newCtx context.Context
