package pgxsql

import (
	"context"
	"errors"
	"github.com/go-ai-agent/core/log"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	upstreamTimeoutFlag = "UT"
	rateLimitedFlag     = "RL"
)

// QueryController - an interface that manages query resiliency
type QueryController interface {
	Apply(ctx context.Context, r Request) (pgx.Rows, *runtime.Status)
}

// ExecController - an interface that manages exec resiliency
type ExecController interface {
	Apply(ctx context.Context, r Request) (pgconn.CommandTag, *runtime.Status)
}

// PingController - an interface that manages ping resiliency
type PingController interface {
	Apply(ctx context.Context) *runtime.Status
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
	limiter   *rate.Limiter
	logFn     startup.AccessLogFn
}

// NewQueryController - create a new resiliency controller
func NewQueryController(name string, threshold Threshold, logFn startup.AccessLogFn) QueryController {
	ctrl := new(controllerCfg)
	ctrl.name = name
	ctrl.threshold = threshold
	if threshold.Limit > 0 {
		ctrl.limiter = rate.NewLimiter(threshold.Limit, threshold.Burst)
	}
	ctrl.logFn = logFn
	return ctrl
}

func (c *controllerCfg) Apply(ctx context.Context, r Request) (rows pgx.Rows, status *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1
	var err error
	logFn := accessFn(ctx, c.logFn)

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		if logFn != nil {
			logFn(log.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, int(c.limiter.Limit()), rateLimitedFlag)
		}
		return
	}
	newCtx := ctx
	if c.threshold.Timeout > 0 {
		childCtx, cancel := context.WithTimeout(ctx, c.threshold.Timeout)
		newCtx = childCtx
		defer cancel()
	}
	if rows, err = dbClient.Query(newCtx, r.Sql(), r.Args()); err != nil {
		if err == context.DeadlineExceeded {
			statusFlags = upstreamTimeoutFlag
			threshold = int(c.threshold.Timeout / time.Millisecond)
			status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
		} else {
			return nil, runtime.NewStatusError(http.StatusInternalServerError, queryLoc, err).SetRequestId(ctx)
		}
	}
	if status == nil {
		status = runtime.NewStatusOK()
	}
	if logFn != nil {
		logFn(log.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, threshold, statusFlags)
	}
	return rows, status
}

type controllerCfgExec controllerCfg

// NewExecController - create a new resiliency controller
func NewExecController(name string, threshold Threshold, logFn startup.AccessLogFn) ExecController {
	ctrl := new(controllerCfgExec)
	ctrl.name = name
	ctrl.threshold = threshold
	if threshold.Limit > 0 {
		ctrl.limiter = rate.NewLimiter(threshold.Limit, threshold.Burst)
	}
	ctrl.logFn = logFn
	return ctrl
}

func (c *controllerCfgExec) Apply(ctx context.Context, r Request) (cmd pgconn.CommandTag, status *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1
	var err error
	logFn := accessFn(ctx, c.logFn)

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return pgconn.CommandTag{}, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		if logFn != nil {
			logFn(log.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, int(c.limiter.Limit()), rateLimitedFlag)
		}
		return
	}
	newCtx := ctx
	if c.threshold.Timeout > 0 {
		childCtx, cancel := context.WithTimeout(ctx, c.threshold.Timeout)
		newCtx = childCtx
		defer cancel()
	}
	if cmd, err = dbClient.Exec(newCtx, r.Sql(), r.Args()); err != nil {
		if err == context.DeadlineExceeded {
			statusFlags = upstreamTimeoutFlag
			threshold = int(c.threshold.Timeout / time.Millisecond)
			status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
		} else {
			return pgconn.CommandTag{}, runtime.NewStatusError(http.StatusInternalServerError, execLoc, err).SetRequestId(ctx)
		}
	}
	if status == nil {
		status = runtime.NewStatusOK()
	}
	if logFn != nil {
		logFn(log.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, threshold, statusFlags)
	}
	return
}

func accessFn(ctx context.Context, logFn startup.AccessLogFn) startup.AccessLogFn {
	if logFn != nil {
		return logFn
	}
	return log.AccessFromContext(ctx)
}

type controllerCfgPing controllerCfg

// NewPingController - create a new resiliency controller
func NewPingController(name string, threshold Threshold, logFn startup.AccessLogFn) PingController {
	ctrl := new(controllerCfgPing)
	ctrl.name = name
	ctrl.threshold = threshold
	ctrl.logFn = logFn
	return ctrl
}

func (c *controllerCfgPing) Apply(ctx context.Context) (status *runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1
	logFn := accessFn(ctx, c.logFn)

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetRequestId(ctx)
	}
	if c.threshold.Timeout > 0 {
		childCtx, cancel := context.WithTimeout(ctx, c.threshold.Timeout)
		defer cancel()
		if err := dbClient.Ping(childCtx); err != nil {
			if err == context.DeadlineExceeded {
				statusFlags = upstreamTimeoutFlag
				threshold = int(c.threshold.Timeout / time.Millisecond)
				status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
			} else {
				return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
			}
		}
	} else {
		if err := dbClient.Ping(ctx); err != nil {
			if err != nil {
				return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
			}
		}
	}
	if status == nil {
		status = runtime.NewStatusOK()
	}
	if logFn != nil {
		req, _ := http.NewRequest(pingControllerName, PingUri, nil)
		resp := http.Response{StatusCode: status.Code()}
		logFn(log.EgressTraffic, start, time.Since(start), req, &resp, threshold, statusFlags)
	}
	return
}