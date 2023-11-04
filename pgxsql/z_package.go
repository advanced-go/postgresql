package pgxsql

import (
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
	"net/http"
	"reflect"
	"sync/atomic"
)

type pkg struct{}

var (
	PkgUri  = reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath = runtime.PathFromUri(PkgUri)
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

// newTypeHandler - templated function providing a TypeHandlerFn via a closure
//func newTypeHandler[E runtime.ErrorHandler]() runtime.TypeHandlerFn {
//	return func(r *http.Request, body any) (any, *runtime.Status) {
//		return typeHandler[E](r, body)
//	}
//}

func TypeHandler(r *http.Request, body any) (any, *runtime.Status) {
	return typeHandler[runtime.LogError](r, body)
}

func typeHandler[E runtime.ErrorHandler](r *http.Request, body any) (any, *runtime.Status) {
	//var e E

	if r == nil {
		return nil, runtime.NewStatus(http.StatusBadRequest)
	}
	if r.URL.Path == startup.StatusPath {
		if isStarted() {
			return nil, runtime.NewStatusOK()
		}
		return nil, runtime.NewStatus(runtime.StatusNotStarted)
	}

	return nil, runtime.NewStatus(http.StatusMethodNotAllowed)
}
