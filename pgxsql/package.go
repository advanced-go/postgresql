package pgxsql

import (
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/startup"
	"net/http"
	"sync/atomic"
)

type pkg struct{}

const (
	PkgUri  = "github.com/advanced-go/postgresql/pgxsql"
	PkgPath = "/advanced-go/postgresql/pgxsql"
)

var (
	started int64
)

// isStarted - returns status of startup
func isStarted() bool {
	return atomic.LoadInt64(&started) != 0
}

func setStarted() {
	atomic.StoreInt64(&started, 1)
}

func resetStarted() {
	atomic.StoreInt64(&started, 0)
}

func GetStartupStatus(uri string) runtime.Status {
	if uri == startup.StatusPath {
		if isStarted() {
			return runtime.NewStatusOK()
		}
		return runtime.NewStatus(runtime.StatusNotStarted)
	}
	return runtime.NewStatus(http.StatusNotFound)
}
