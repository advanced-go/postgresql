package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
)

const (
	execLoc = PkgPath + ":exec"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status runtime.Status) {
	url, override := lookup.Value(req.resource)
	var newCtx context.Context

	if req == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetRequestId(ctx)
	}
	defer apply(ctx, &newCtx, req, statusCode(&status))
	if override {
		return runtime.New[CommandTag](url, nil)
	}
	if dbClient == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetRequestId(ctx)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(newCtx)
	if err0 != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err0).SetRequestId(newCtx)
	}
	cmd, err := dbClient.Exec(newCtx, buildSql(req), req.args)
	if err != nil {
		err0 = txn.Rollback(newCtx)
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, recast(err), err0).SetRequestId(newCtx)
	}
	err = txn.Commit(newCtx)
	if err != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err).SetRequestId(newCtx)
	}
	return newCmdTag(cmd), runtime.StatusOK()
}
