package pgxsql

import (
	"context"
	"github.com/advanced-go/core/runtime"
)

type statusT struct{}

var (
	statusKey = statusT{}
)

// NewStatusContext - creates a new Context with a Status
func NewStatusContext(ctx context.Context, status runtime.Status) context.Context {
	if ctx == nil {
		ctx = context.Background()
	} else {
		i := ctx.Value(statusKey)
		if i != nil {
			return ctx
		}
	}
	return context.WithValue(ctx, statusKey, status)
}

// StatusFromContext - return a Status from a context
func StatusFromContext(ctx any) runtime.Status {
	if ctx == nil {
		return nil
	}
	if ctx2, ok := ctx.(context.Context); ok {
		i := ctx2.Value(statusKey)
		if status, ok2 := i.(runtime.Status); ok2 {
			return status
		}
	}
	return nil
}
