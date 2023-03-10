package pgxsql

import (
	"context"
	"time"
)

// ExecExchange - interface for exec
type ExecExchange interface {
	Exec(*Request) (CommandTag, error)
}

func execExchangeCast(ctx context.Context) (ExecExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(ExecExchange); ok {
		return e, true
	}
	return nil, false
}

//type ExecContext interface {
//	context.Context
//	ExecExchange
//	withValue(key, val any) context.Context
//}

type execWithValue interface {
	execWithValue(key, val any) context.Context
}

type execContext struct {
	ctx  context.Context
	exec func(*Request) (CommandTag, error)
}

func ContextWithExec(ctx context.Context, exec func(*Request) (CommandTag, error)) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &execContext{ctx: ctx, exec: exec}
}

func (c *execContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *execContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *execContext) Err() error {
	return c.ctx.Err()
}

func (c *execContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *execContext) Exec(req *Request) (CommandTag, error) {
	return c.exec(req)
}

func (c *execContext) execWithValue(key, val any) context.Context {
	c.ctx = context.WithValue(c.ctx, key, val)
	return c
}
