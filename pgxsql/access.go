package pgxsql

import (
	"github.com/advanced-go/stdlib/access"
	"net/http"
	"time"
)

func log(start time.Time, duration time.Duration, req request, routeName string, threshold int, flags string) {
	var r *http.Request
	var resp *http.Response

	access.Log(access.EgressTraffic, start, duration, r, resp, routeName, "", threshold, flags)
}
