package pgxsql

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

type pkg struct{}

const (
	PkgPath         = "github.com/advanced-go/postgresql/pgxsql"
	ReadinessPath   = PkgPath + ":Readiness"
	ContentLocation = "Content-Location"
)

func Readiness(uri string) runtime.Status {
	if uri == ReadinessPath {
		if isReady() {
			return runtime.StatusOK()
		}
		return runtime.NewStatus(runtime.StatusNotStarted)
	}
	return runtime.NewStatus(http.StatusNotFound)
}
