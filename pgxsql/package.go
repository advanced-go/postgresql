package pgxsql

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
	"sync/atomic"
)

type pkg struct{}

const (
	PkgPath         = "github.com/advanced-go/postgresql/pgxsql"
	StatusPath      = PkgPath + ":Status"
	ContentLocation = "Content-Location"
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
	if uri == StatusPath {
		if isStarted() {
			return runtime.StatusOK()
		}
		return runtime.NewStatus(runtime.StatusNotStarted)
	}
	return runtime.NewStatus(http.StatusNotFound)
}
