package pgxsql

import (
	"context"
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
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
	Apply(ctx context.Context, r Request) (pgx.Rows, runtime.Status)
	getThreshold() Threshold
	updateRateLimiter(limit rate.Limit, burst int)
}

// ExecController - an interface that manages exec resiliency
type ExecController interface {
	Apply(ctx context.Context, r Request) (pgconn.CommandTag, runtime.Status)
	getThreshold() Threshold
	updateRateLimiter(limit rate.Limit, burst int)
}

// PingController - an interface that manages ping resiliency
type PingController interface {
	Apply(ctx context.Context) runtime.Status
	getThreshold() Threshold
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
	logFn     access.LogHandler
}

// NewQueryController - create a new resiliency controller
func NewQueryController(name string, threshold Threshold, logFn access.LogHandler) QueryController {
	ctrl := new(controllerCfg)
	ctrl.name = name
	ctrl.threshold = threshold
	if threshold.Limit > 0 {
		ctrl.limiter = rate.NewLimiter(threshold.Limit, threshold.Burst)
	}
	ctrl.logFn = logFn
	if ctrl.logFn == nil {
		ctrl.logFn = access.Log
	}
	return ctrl
}

func (c *controllerCfg) Apply(ctx context.Context, r Request) (rows pgx.Rows, status runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1
	var err error

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, queryLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, int(c.limiter.Limit()), rateLimitedFlag)
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
		status = runtime.StatusOK()
	}
	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, threshold, statusFlags)
	return rows, status
}

func (c *controllerCfg) getThreshold() Threshold {
	return c.threshold
}

func (c *controllerCfg) updateRateLimiter(limit rate.Limit, burst int) {
	c.limiter.SetLimit(limit)
	c.limiter.SetBurst(burst)
}

type controllerCfgExec controllerCfg

// NewExecController - create a new resiliency controller
func NewExecController(name string, threshold Threshold, logFn access.LogHandler) ExecController {
	ctrl := new(controllerCfgExec)
	ctrl.name = name
	ctrl.threshold = threshold
	if threshold.Limit > 0 {
		ctrl.limiter = rate.NewLimiter(threshold.Limit, threshold.Burst)
	}
	ctrl.logFn = logFn
	if ctrl.logFn == nil {
		ctrl.logFn = access.Log
	}
	return ctrl
}

func (c *controllerCfgExec) Apply(ctx context.Context, r Request) (cmd pgconn.CommandTag, status runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1
	var err error

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return pgconn.CommandTag{}, runtime.NewStatusError(runtime.StatusInvalidArgument, execLoc, errors.New("error on PostgreSQL database query call: dbClient is nil")).SetRequestId(ctx)
	}
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, int(c.limiter.Limit()), rateLimitedFlag)
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
		status = runtime.StatusOK()
	}
	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, threshold, statusFlags)
	return
}

func (c *controllerCfgExec) getThreshold() Threshold {
	return c.threshold
}

func (c *controllerCfgExec) updateRateLimiter(limit rate.Limit, burst int) {
	c.limiter.SetLimit(limit)
	c.limiter.SetBurst(burst)
}

type controllerCfgPing controllerCfg

// NewPingController - create a new resiliency controller
func NewPingController(name string, threshold Threshold, logFn access.LogHandler) PingController {
	ctrl := new(controllerCfgPing)
	ctrl.name = name
	ctrl.threshold = threshold
	ctrl.logFn = logFn
	if ctrl.logFn == nil {
		ctrl.logFn = access.Log
	}
	return ctrl
}

func (c *controllerCfgPing) Apply(ctx context.Context) (status runtime.Status) {
	start := time.Now().UTC()
	statusFlags := ""
	threshold := -1

	if ctx == nil {
		ctx = context.Background()
	}
	if dbClient == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, pingLoc, errors.New("error on PostgreSQL ping call : dbClient is nil")).SetRequestId(ctx)
	}
	//var newCtx context.Context
	if c.threshold.Timeout > 0 {
		newCtx, cancel := context.WithTimeout(ctx, c.threshold.Timeout)
		ctx = newCtx
		defer cancel()
	}
	err := dbClient.Ping(ctx)
	dur := time.Since(start)
	if err != nil {
		if err == context.DeadlineExceeded {
			statusFlags = upstreamTimeoutFlag
			threshold = int(c.threshold.Timeout / time.Millisecond)
			status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
		} else {
			return runtime.NewStatusError(http.StatusInternalServerError, pingLoc, err).SetRequestId(ctx)
		}
	}
	if status == nil {
		status = runtime.StatusOK()
	}
	//status.SetDuration(dur)
	req, _ := http.NewRequest(pingControllerName, PingUri, nil)
	c.logFn(access.EgressTraffic, start, dur, req, &http.Response{StatusCode: status.Code()}, threshold, statusFlags)
	return
}

func (c *controllerCfgPing) getThreshold() Threshold {
	return c.threshold
}
