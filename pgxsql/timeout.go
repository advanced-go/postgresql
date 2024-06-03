package pgxsql

import (
	"context"
	"time"
)

func setTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, duration)
}
