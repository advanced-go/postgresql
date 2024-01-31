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
func apply2(ctx context.Context, r *request, status *runtime.Status) (func(), context.Context) {
	thresholdFlags := ""
	start := time.Now()
	newCtx := ctx
	var cancelFunc context.CancelFunc
	req, _ := http.NewRequest(method(r), r.uri, nil)
	if r.header != nil {
		req.Header = r.header
	}

	// TO DO : determine if the current context already contains a CancelCtx
	if ctx != nil {
	} else {
		newCtx, cancelFunc = context.WithTimeout(context.Background(), time.Millisecond*time.Duration(r.threshold))
	}
	return func() {
		if cancelFunc != nil {
			cancelFunc()
		}
		threshold := r.threshold
		code := (*status).Code()
		if code == runtime.StatusDeadlineExceeded {
			thresholdFlags = upstreamTimeoutFlag
		} else {
			threshold = -1
		}
		access.Log(access.EgressTraffic, start, time.Since(start), req, &http.Response{StatusCode: code, Status: (*status).Description()}, r.routeName, "", threshold, thresholdFlags)
	}, newCtx
}
