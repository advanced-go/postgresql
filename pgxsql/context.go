package pgxsql

import "context"

// ContextWithValue - create a new context with a value, updating the context if it is an ExecExchange or QueryExchange context
func ContextWithValue(ctx context.Context, key any, val any) context.Context {
	if ctx == nil {
		return nil
	}
	if curr, ok := any(ctx).(execWithValue); ok {
		return curr.execWithValue(key, val)
	}
	if curr, ok := any(ctx).(queryWithValue); ok {
		return curr.queryWithValue(key, val)
	}
	return context.WithValue(ctx, key, val)
}
