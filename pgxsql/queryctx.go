package pgxsql

import (
	"context"
	"time"
)

// QueryExchange - interface for query
type QueryExchange interface {
	Query(req *Request) (Rows, error)
}

type queryWithValue interface {
	queryWithValue(key, val any) context.Context
}

type queryContext struct {
	ctx   context.Context
	query func(req *Request) (Rows, error)
}

func ContextWithQuery(ctx context.Context, query func(req *Request) (Rows, error)) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &queryContext{ctx: ctx, query: query}
}

func (c *queryContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *queryContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *queryContext) Err() error {
	return c.ctx.Err()
}

func (c *queryContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *queryContext) Query(req *Request) (Rows, error) {
	return c.query(req)
}

func (c *queryContext) queryWithValue(key, val any) context.Context {
	c.ctx = context.WithValue(c.ctx, key, val)
	return c
}

func queryExchangeCast(ctx context.Context) (QueryExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(QueryExchange); ok {
		return e, true
	}
	return nil, false
}
