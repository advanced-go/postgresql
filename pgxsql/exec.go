package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sre/core/runtime"
	"github.com/go-sre/core/sql"
	"github.com/go-sre/host/messaging"
)

const (
	NullCount = int64(-1)
)

var execLoc = pkgPath + "/exec"

// Exec - templated function for executing a SQL statement
func Exec[E runtime.ErrorHandler](ctx context.Context, req *sql.Request, args ...any) (tag CommandTag, status *runtime.Status) {
	var e E
	var limited = false
	var fn func()

	if ctx == nil {
		ctx = context.Background()
	}
	if req == nil {
		return tag, e.Handle(ctx, execLoc, errors.New("error on PostgreSQL exec call : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	fn, ctx, limited = controllerApply(ctx, messaging.NewStatusCode(&status), req.Uri, runtime.ContextRequestId(ctx), "GET")
	defer fn()
	if limited {
		return tag, runtime.NewStatusCode(runtime.StatusRateLimited)
	}
	//if exec, ok := execExchangeCast(ctx); ok {
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		if pExec := findExecProxy(proxies); pExec != nil {
			result, err := pExec(req)
			return result, e.Handle(ctx, execLoc, err)
		}
	}
	if dbClient == nil {
		return tag, e.Handle(ctx, execLoc, errors.New("error on PostgreSQL exec call : dbClient is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(ctx)
	if err0 != nil {
		return tag, e.Handle(ctx, execLoc, err0)
	}
	t, err := dbClient.Exec(ctx, BuildSql(req), args...)
	if err != nil {
		err0 = txn.Rollback(ctx)
		return tag, e.Handle(ctx, execLoc, recast(err), err0)
	}
	if req.ExpectedCount != sql.NullExpectedCount && t.RowsAffected() != req.ExpectedCount {
		err0 = txn.Rollback(ctx)
		return tag, e.Handle(ctx, execLoc, errors.New(fmt.Sprintf("error exec statement [%v] : actual RowsAffected %v != expected RowsAffected %v", t.String(), t.RowsAffected(), req.ExpectedCount)), err0)
	}
	err = txn.Commit(ctx)
	if err != nil {
		return tag, e.Handle(ctx, execLoc, err)
	}
	return CommandTag{Sql: t.String(), RowsAffected: t.RowsAffected(), Insert: t.Insert(), Update: t.Update(), Delete: t.Delete(), Select: t.Select()}, runtime.NewStatusOK()
}
