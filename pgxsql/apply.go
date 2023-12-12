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

// Apply - function to be used to access log and apply a timeout
func apply(ctx context.Context, r Request, routeName string, threshold int, statusCode func() int) (func(), context.Context) {
	thresholdFlags := ""
	start := time.Now()
	newCtx := ctx
	var cancelFunc context.CancelFunc
	req, _ := http.NewRequest(r.Method(), r.Uri(), nil)
	if r.Header() != nil {
		req.Header = r.Header()
	}

	// TO DO : determine if the current context already contains a CancelCtx
	if ctx != nil {
	} else {
		newCtx, cancelFunc = context.WithTimeout(context.Background(), time.Millisecond*time.Duration(threshold))
	}
	return func() {
		if cancelFunc != nil {
			cancelFunc()
		}
		code := statusCode()
		if code == runtime.StatusDeadlineExceeded {
			thresholdFlags = upstreamTimeoutFlag
		} else {
			threshold = -1
		}
		access.Log(access.EgressTraffic, start, time.Since(start), req, &http.Response{StatusCode: code}, routeName, threshold, thresholdFlags)
	}, newCtx
}
