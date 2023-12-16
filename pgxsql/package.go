package pgxsql

import (
	"context"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type pkg struct{}

const (
	PkgPath = "github.com/advanced-go/postgresql/pgxsql"
)

// Credentials - credentials function for authentication
type Credentials func() (username string, password string, err error)

// Resource - database URL for connectivity configuration
type Resource struct {
	Uri string
}

// Attr - key value pair
type Attr struct {
	Key string
	Val any
}

// Readiness - package readiness
func Readiness() runtime.Status {
	if isReady() {
		return runtime.StatusOK()
	}
	return runtime.NewStatus(runtime.StatusNotStarted)
}

// Query -  process a SQL select statement
func Query(ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows pgx.Rows, status runtime.Status) {
	req := newQueryRequestFromValues(h, resource, template, values, args...)
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(h, method(req), req.uri), req.routeName, -1, "", access.NewStatusCodeClosure(&status))()
	return query(ctx, req)
}

// Insert - execute a SQL insert statement
func Insert(ctx context.Context, h http.Header, resource, template string, values [][]any, args ...any) (tag CommandTag, status runtime.Status) {
	req := newInsertRequest(h, resource, template, values, args...)
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(h, method(req), req.uri), req.routeName, -1, "", access.NewStatusCodeClosure(&status))()
	return exec(ctx, req)
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, template string, where []Attr, args []Attr) (tag CommandTag, status runtime.Status) {
	req := newUpdateRequest(h, resource, template, convert(where), convert(args))
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(h, method(req), req.uri), req.routeName, -1, "", access.NewStatusCodeClosure(&status))()
	return exec(ctx, req)
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, template string, where []Attr, args ...any) (tag CommandTag, status runtime.Status) {
	req := newDeleteRequest(h, resource, template, convert(where), args...)
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(h, method(req), req.uri), req.routeName, -1, "", access.NewStatusCodeClosure(&status))()
	return exec(ctx, req)
}

// Stat - retrieve runtime stats
func Stat(ctx context.Context) (*pgxpool.Stat, runtime.Status) {
	return stat(ctx)
}

// Ping - ping the database cluster
func Ping(ctx context.Context) runtime.Status {
	return ping(ctx)
}
