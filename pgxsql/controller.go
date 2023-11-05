package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/log"
	"github.com/go-ai-agent/core/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	egressTraffic = "egress"
)

// QueryController - an interface that manages query resiliency
type QueryController interface {
	Apply(ctx context.Context, r Request) (pgx.Rows, *runtime.Status)
}

// ExecController - an interface that manages exec resiliency
type ExecController interface {
	Apply(ctx context.Context, r Request) (pgconn.CommandTag, *runtime.Status)
}

// Threshold - rate limiting and timeout
type Threshold struct {
	Limit   rate.Limit // request per second
	Burst   int
	Timeout time.Duration
}

type controllerCfg struct {
	name      string
	threshold Threshold
}

// NewQueryController - create a new resiliency controller
func NewQueryController(name string, threshold Threshold) QueryController {
	ctrl := new(controllerCfg)
	ctrl.name = name
	ctrl.threshold = threshold
	return ctrl
}

func (c *controllerCfg) Apply(ctx context.Context, r Request) (pgx.Rows, *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""

	rows, status := applyQuery(ctx, r)
	logger := log.ContextAccessLogger(ctx)
	if logger != nil {
		req, _ := http.NewRequest(r.Method(), r.Uri(), nil)
		resp := http.Response{StatusCode: status.Code()}
		logger(egressTraffic, start, time.Since(start), req, &resp, statusFlags) // c.name, c.threshold.Limit, c.threshold.Burst, int(c.threshold.Timeout/time.Millisecond), statusFlags)
	}
	return rows, status
}

type controllerCfgExec controllerCfg

// NewExecController - create a new resiliency controller
func NewExecController(name string, threshold Threshold) ExecController {
	ctrl := new(controllerCfgExec)
	ctrl.name = name
	ctrl.threshold = threshold
	return ctrl
}

func (c *controllerCfgExec) Apply(ctx context.Context, r Request) (pgconn.CommandTag, *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""

	cmd, status := applyExec(ctx, r)
	logger := log.ContextAccessLogger(ctx)
	if logger != nil {
		req, _ := http.NewRequest(r.Method(), r.Uri(), nil)
		resp := http.Response{StatusCode: status.Code()}
		logger(egressTraffic, start, time.Since(start), req, &resp, statusFlags) // c.name, c.threshold.Limit, c.threshold.Burst, int(c.threshold.Timeout/time.Millisecond), statusFlags)
	}
	return cmd, status
}

func applyQuery(ctx context.Context, r Request) (pgx.Rows, *runtime.Status) {
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	rows, err := dbClient.Query(ctx, r.Sql(), r.Args())
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, err)
	}
	return rows, runtime.NewStatusOK()
}

func applyExec(ctx context.Context, r Request) (pgconn.CommandTag, *runtime.Status) {
	if dbClient == nil {
		return pgconn.CommandTag{}, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	cmd, err := dbClient.Exec(ctx, r.Sql(), r.Args())
	if err != nil {
		return pgconn.CommandTag{}, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, err)
	}
	return cmd, runtime.NewStatusOK()
}
