package pgxsql

import (
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/runtime/startup"
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

func GetStatus() *runtime.Status {
	_, status := doHandler[runtime.LogError](nil, "", startup.StatusPath, "", nil)
	return status
}

func doHandler[E runtime.ErrorHandler](_ any, _, uri, _ string, _ any) (any, *runtime.Status) {
	if uri == startup.StatusPath {
		if isStarted() {
			return nil, runtime.NewStatusOK()
		}
		return nil, runtime.NewStatus(runtime.StatusNotStarted)
	}
	return nil, runtime.NewStatus(http.StatusMethodNotAllowed)
}
