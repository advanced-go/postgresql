package pgxsql

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type queryBypass controllerCfg

func NewQueryBypassController(name string, log AccessLogFn) QueryController {
	ctrl := new(queryBypass)
	ctrl.name = name
	ctrl.log = log
	return ctrl
}

func (b *queryBypass) Apply(ctx context.Context, r Request) (pgx.Rows, *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""

	rows, status := applyQuery(ctx, r)
	if b.log != nil {
		b.log(egressTraffic, start, time.Since(start), r.Uri(), r.Method(), status.Code(), b.name, b.threshold.Limit, b.threshold.Burst, int(b.threshold.Timeout/time.Millisecond), statusFlags)
	}
	return rows, status

}

type execBypass controllerCfg

func NewExecBypassController(name string, log AccessLogFn) ExecController {
	ctrl := new(execBypass)
	ctrl.name = name
	ctrl.log = log
	return ctrl
}

func (b *execBypass) Apply(ctx context.Context, r Request) (pgconn.CommandTag, *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""

	cmd, status := applyExec(ctx, r)
	if b.log != nil {
		b.log(egressTraffic, start, time.Since(start), r.Uri(), r.Method(), status.Code(), b.name, b.threshold.Limit, b.threshold.Burst, int(b.threshold.Timeout/time.Millisecond), statusFlags)
	}
	return cmd, status
}
