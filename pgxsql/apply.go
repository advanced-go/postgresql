package pgxsql

import (
	"context"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	upstreamTimeoutFlag = "UT"
)

// Apply - function to be used to apply a controller
func Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	statusFlags := ""
	limited := false
	start := time.Now()
	newCtx := ctx
	var cancelCtx context.CancelFunc
	req, _ := http.NewRequest(method, uri, nil)
	threshold := 500

	/*
		ctrl, err := EgressLookup(req)
		if err != nil {
			statusFlags = err.Error()
		} else {
			if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
				limited = true
				statusFlags = RateLimitFlag
			}
			if !limited {
				if to := ctrl.Timeout(); to.IsEnabled() {
					newCtx, cancelCtx = context.WithTimeout(ctx, to.Duration())
				}
			}
		}

	*/

	////if to := ctrl.Timeout(); to.IsEnabled() {
	//	newCtx, cancelCtx = context.WithTimeout(ctx, to.Duration())
	//}
	return func() {
		if cancelCtx != nil {
			cancelCtx()
		}
		code := statusCode()
		if code == runtime.StatusDeadlineExceeded {
			statusFlags = upstreamTimeoutFlag
		} else {
			threshold = -1
		}
		access.Log(access.EgressTraffic, start, time.Since(start), req, &http.Response{StatusCode: code}, "", threshold, statusFlags)
	}, newCtx, limited
}
