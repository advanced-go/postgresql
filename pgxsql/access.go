package pgxsql

import (
	"github.com/advanced-go/postgresql/module"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

func log(start time.Time, duration time.Duration, req *request, status *core.Status, flags string) {
	r := NewHttpRequest(req)
	resp := &http.Response{StatusCode: status.HttpCode()}
	r.Header.Add(core.XAuthority, module.Authority)

	access.Log(access.EgressTraffic, start, duration, r, resp, req.routeName, "", milliseconds(req.duration), flags)
}

func milliseconds(duration time.Duration) int {
	return int(duration / time.Duration(1e6))
}
