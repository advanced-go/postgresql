package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status *core.Status) {
	if req == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : request is nil"))
	}
	if dbClient == nil {
		status = core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : dbClient is nil"))
		return
	}
	ctx = req.setTimeout(ctx)

	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		status = core.NewStatusError(core.StatusTxnBeginError, err0)
		return tag, status
	}
	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer txn.Rollback(ctx)
	cmd, err := dbClient.Exec(ctx, buildSql(req), req.args)
	if err != nil {
		status = core.NewStatusError(core.StatusInvalidArgument, recast(err))
		return newCmdTag(cmd), status
	}
	err = txn.Commit(ctx)
	if err != nil {
		status = core.NewStatusError(core.StatusTxnCommitError, err)
	} else {
		status = core.StatusOK()
	}
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
