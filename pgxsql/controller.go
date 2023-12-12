package pgxsql

import (
	"context"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const (
	//upstreamTimeoutFlag = "UT"
	rateLimitedFlag = "RL"
)

// QueryController - an interface that manages query resiliency
type queryControllerT interface {
	Apply(ctx context.Context, r Request) (pgx.Rows, runtime.Status)
	getThreshold() thresholdValues
	updateRateLimiter(limit rate.Limit, burst int)
}

// ExecController - an interface that manages exec resiliency
type execControllerT interface {
	Apply(ctx context.Context, r Request) (pgconn.CommandTag, runtime.Status)
	getThreshold() thresholdValues
	updateRateLimiter(limit rate.Limit, burst int)
}

// PingController - an interface that manages ping resiliency
type pingControllerT interface {
	Apply(ctx context.Context) runtime.Status
	getThreshold() thresholdValues
}

// Threshold - rate limiting and timeout
type thresholdValues struct {
	Limit   rate.Limit // request per second
	Burst   int
	Timeout time.Duration
}

type controllerCfg struct {
	name      string
	threshold thresholdValues
	limiter   *rate.Limiter
	logFn     access.LogHandler
}

// newQueryController - create a new resiliency controller
func newQueryController(name string, threshold thresholdValues, logFn access.LogHandler) queryControllerT {
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
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, c.name, int(c.limiter.Limit()), rateLimitedFlag)
		return
	}
	// Override
	//location := r.Header().Get(ContentLocation)
	//if len(location) > 0 {
	// TO DO : Need to read the Rows from the location
	//	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: http.StatusOK}, c.name, threshold, statusFlags)
	//	return nil, runtime.StatusOK()
	//}
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
	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, c.name, threshold, statusFlags)
	return rows, status
}

func (c *controllerCfg) getThreshold() thresholdValues {
	return c.threshold
}

func (c *controllerCfg) updateRateLimiter(limit rate.Limit, burst int) {
	c.limiter.SetLimit(limit)
	c.limiter.SetBurst(burst)
}

type controllerCfgExec controllerCfg

// newExecController - create a new resiliency controller
func newExecController(name string, threshold thresholdValues, logFn access.LogHandler) execControllerT {
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
	if c.limiter != nil && !c.limiter.Allow() {
		status = runtime.NewStatus(runtime.StatusRateLimited)
		c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, c.name, int(c.limiter.Limit()), rateLimitedFlag)
		return
	}
	// Override
	//location := r.Header().Get(ContentLocation)
	//if len(location) > 0 {
	// TO DO : Need to read the command tag from the location
	//	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: http.StatusOK}, c.name, threshold, statusFlags)
	//	return pgconn.CommandTag{}, runtime.StatusOK()
	//}
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
	c.logFn(access.EgressTraffic, start, time.Since(start), r.HttpRequest(), &http.Response{StatusCode: status.Code()}, c.name, threshold, statusFlags)
	return
}

func (c *controllerCfgExec) getThreshold() thresholdValues {
	return c.threshold
}

func (c *controllerCfgExec) updateRateLimiter(limit rate.Limit, burst int) {
	c.limiter.SetLimit(limit)
	c.limiter.SetBurst(burst)
}

type controllerCfgPing controllerCfg

// newPingController - create a new resiliency controller
func newPingController(name string, threshold thresholdValues, logFn access.LogHandler) pingControllerT {
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
	c.logFn(access.EgressTraffic, start, dur, req, &http.Response{StatusCode: status.Code()}, c.name, threshold, statusFlags)
	return
}

func (c *controllerCfgPing) getThreshold() thresholdValues {
	return c.threshold
}
