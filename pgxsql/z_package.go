package pgxsql

import (
	"github.com/go-ai-agent/core/runtime"
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

// IsStarted - returns status of startup
func IsStarted() bool {
	return atomic.LoadInt64(&started) != 0
}

func setStarted() {
	atomic.StoreInt64(&started, 1)
}

func resetStarted() {
	atomic.StoreInt64(&started, 0)
}

// newTypeHandler - templated function providing a TypeHandlerFn via a closure
func newTypeHandler[E runtime.ErrorHandler]() runtime.TypeHandlerFn {
	return func(r *http.Request, body any) (any, *runtime.Status) {
		return typeHandler[E](r, body)
	}
}

func typeHandler[E runtime.ErrorHandler](r *http.Request, body any) (any, *runtime.Status) {
	//var e E

	if r == nil {
		return nil, runtime.NewStatus(http.StatusBadRequest)
	}
	// create a new context with a request id. Not creating a new request as upstream processing doesn't
	// use http
	/*
		requestId := runtime.GetOrCreateRequestId(r)
		nc := runtime.ContextWithRequestId(r.Context(), requestId)
		switch r.Method {
		case http.MethodGet:
			entries, status := get(nc, r.Header.Get(httpx.ContentLocation), r.URL.Query())
			if !status.OK() {
				e.HandleStatus(status, requestId, locTypeHandler)
				return nil, status
			}
			if entries == nil {
				status.SetCode(http.StatusNotFound)
			}
			return entries, status
		case http.MethodPut:
			cmdTag, status := put(nc, r.Header.Get(httpx.ContentLocation), body)
			if !status.OK() {
				e.HandleStatus(status, requestId, locTypeHandler)
				return nil, status
			}
			return cmdTag, status
		default:
		}

	*/
	return nil, runtime.NewStatus(http.StatusMethodNotAllowed)
}
