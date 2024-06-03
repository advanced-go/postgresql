package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
)

const (
	execLoc = PkgPath + ":exec"
)

func exec(ctx context.Context, req *request) (tag CommandTag, status *core.Status) {
	//url, override := lookup.Value(req.resource)
	var newCtx context.Context

	if req == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : request is nil"))
	}
	//defer apply(ctx, &newCtx, req, access.StatusCode(&status))
	//if override {
	//	return io2.New[CommandTag](url, nil)
	//}
	if dbClient == nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL exec call : dbClient is nil"))
	}
	// Transaction processing.
	txn, err0 := dbClient.Begin(newCtx)
	if err0 != nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, err0)
	}
	cmd, err := dbClient.Exec(newCtx, buildSql(req), req.args)
	if err != nil {
		err0 = txn.Rollback(newCtx)
		// TODO: how to handle err0?
		return tag, core.NewStatusError(core.StatusInvalidArgument, recast(err))
	}
	err = txn.Commit(newCtx)
	if err != nil {
		return tag, core.NewStatusError(core.StatusInvalidArgument, err)
	}
	return newCmdTag(cmd), core.StatusOK()
}
