package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context, req *request) (status *core.Status) {
	if dbClient == nil {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	ctx1, cancel := setTimeout(ctx, req)
	if cancel != nil {
		defer cancel()
	}
	var start = time.Now().UTC()
	err := dbClient.Ping(ctx1)
	if err != nil {
		status = core.NewStatusError(http.StatusInternalServerError, err)
	} else {
		status = core.StatusOK()
	}
	// TODO : determine if there was a timeout
	log(start, time.Since(start), req, status, "")
	return
}

// Scrap
//var limited = false
//var fn func()
//
