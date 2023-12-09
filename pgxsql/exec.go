package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
)

const (
	execLoc = PkgPath + ":Exec"
)

// Exec - function for executing a SQL statement
func Exec(ctx context.Context, req Request) (tag CommandTag, status runtime.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetRequestId(ctx)
	}
	if runtime.IsDebugEnvironment() {
		status = StatusFromContext(ctx)
		if status != nil {
			return CommandTag{}, status
		}
		if r, ok := any(req).(*request); ok {
			if r.execProxy() != nil {
				return r.execProxy()(req)
			}
		}
	}
	if dbClient == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetRequestId(ctx)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err0).SetRequestId(ctx)
	}
	t, status1 := execController.Apply(ctx, req)
	if !status1.OK() {
		err0 = txn.Rollback(ctx)
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, recast(status1.FirstError()), err0).SetRequestId(ctx)
	}
	err := txn.Commit(ctx)
	if err != nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err).SetRequestId(ctx)
	}
	return newCmdTag(t), status1
}

//t, err := dbClient.Exec(ctx, req.Sql(), req.Args())
