package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/runtime"
)

var execLoc = PkgUri + "/Exec"

// Exec - function for executing a SQL statement
func Exec(ctx context.Context, req Request) (tag CommandTag, status *runtime.Status) {
	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return tag, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetRequestId(ctx)
	}
	if runtime.IsDebugEnvironment() {
		if proxies, ok := runtime.IsProxyable(ctx); ok {
			if pExec := findExecProxy(proxies); pExec != nil {
				result, err := pExec(req)
				return newCmdTag(result), runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err).SetRequestId(ctx)
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
