package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
)

const (
	execLoc = PkgPath + ":exec"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status runtime.Status) {
	var fn func()
	urls, ok := lookup(req.resource)

	if req == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetRequestId(ctx)
	}
	if !ok && dbClient == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetRequestId(ctx)
	}
	fn, ctx = apply(ctx, req, &status)
	defer fn()
	if ok {
		return io2.ReadResults[CommandTag](urls)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err0).SetRequestId(ctx)
	}
	cmd, err := dbClient.Exec(ctx, buildSql(req), req.args)
	if err != nil {
		err0 = txn.Rollback(ctx)
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, recast(err), err0).SetRequestId(ctx)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err).SetRequestId(ctx)
	}
	return newCmdTag(cmd), runtime.StatusOK()
}
