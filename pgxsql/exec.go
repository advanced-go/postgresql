package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
)

const (
	execLoc       = PkgPath + ":exec"
	execRouteName = "exec"
)

func exec(ctx context.Context, req Request) (tag CommandTag, status runtime.Status) {
	var fn func()

	if req == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetRequestId(ctx)
	}
	if dbClient == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetRequestId(ctx)
	}
	fn, ctx = apply(ctx, access.NewStatusCodeClosure(&status), req.Uri(), runtime.RequestId(ctx), req.Method(), execRouteName, execThreshold)
	defer fn()
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err0).SetRequestId(ctx)
	}
	cmd, err := dbClient.Exec(ctx, req.Sql(), req.Args())
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
