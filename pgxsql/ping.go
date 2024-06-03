package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

// Ping - function for pinging the database cluster
func ping(ctx context.Context, duration time.Duration) (status *core.Status) {
	if dbClient == nil {
		return core.NewStatusError(core.StatusInvalidArgument, errors.New("error on PostgreSQL ping call : dbClient is nil"))
	}
	ctx1, cancel := setTimeout(ctx, duration)
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
	logPing(start, time.Since(start), status, access.Milliseconds(duration), "")
	return
}

// Scrap
//var limited = false
//var fn func()
//
