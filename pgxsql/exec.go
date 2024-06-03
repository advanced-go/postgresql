package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
	"time"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status *core.Status) {
	if req == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : request is nil"))
	}

	if dbClient == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : dbClient is nil"))
	}
	ctx1, cancel := setTimeout(ctx, req)
	if cancel != nil {
		defer cancel()
	}
	var start = time.Now().UTC()
	if req.execFunc != nil {
		cmd, err := req.execFunc(ctx1, buildSql(req), req.args)
		status = core.NewStatusError(core.StatusInvalidArgument, recast(err))
		// TODO : determine if there was a timeout
		log(start, time.Since(start), req, status, "")
		return newCmdTag(cmd), status
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx1)
	if err0 != nil {
		status = core.NewStatusError(core.StatusTxnBeginError, err0)
		// TODO : determine if there was a timeout
		log(start, time.Since(start), req, status, "")
		return tag, status
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx1)
	cmd, err := dbClient.Exec(ctx1, buildSql(req), req.args)
	if err != nil {
		status = core.NewStatusError(core.StatusInvalidArgument, recast(err))
		// TODO : determine if there was a timeout
		log(start, time.Since(start), req, status, "")
		return newCmdTag(cmd), status
	}
	err = txn.Commit(ctx1)
	if err != nil {
		status = core.NewStatusError(core.StatusTxnCommitError, err)
	} else {
		status = core.StatusOK()
	}
	// TODO : determine if there was a timeout
	log(start, time.Since(start), req, status, "")
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
