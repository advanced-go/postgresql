package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/common/core"
	"net/http"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context, req *request) (status *core.Status) {
	if dbClient == nil {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	ctx = req.setTimeout(ctx)
	err := dbClient.Ping(ctx)
	if err != nil {
		status = core.NewStatusError(http.StatusInternalServerError, err)
	} else {
		status = core.StatusOK()
	}
	return
}
