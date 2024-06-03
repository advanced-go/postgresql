package pgxsql

import (
	"github.com/advanced-go/postgresql/module"
	"github.com/advanced-go/stdlib/access"
	"github.com/advanced-go/stdlib/core"
	"net/http"
	"time"
)

func log(start time.Time, duration time.Duration, req *request, status *core.Status, threshold int, flags string) {
	r := NewHttpRequest(req)
	resp := &http.Response{StatusCode: status.HttpCode()}

	access.Log(access.EgressTraffic, start, duration, r, resp, req.routeName, "", threshold, flags)
}

func logPing(start time.Time, duration time.Duration, status *core.Status, threshold int, flags string) {
	r, _ := http.NewRequest(pingResource, module.Authority+":"+pingResource, nil)
	resp := &http.Response{StatusCode: status.HttpCode()}

	access.Log(access.EgressTraffic, start, duration, r, resp, PingRouteName, "", threshold, flags)
}
