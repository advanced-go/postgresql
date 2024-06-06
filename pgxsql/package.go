package pgxsql

import (
	"context"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

type pkg struct{}

const (
	PkgPath       = "github/advanced-go/postgresql/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"

	// Timeouts
	QueryTimeout  = time.Second * 2
	InsertTimeout = time.Second * 2
	UpdateTimeout = time.Second * 2
	DeleteTimeout = time.Second * 2
	PingTimeout   = time.Second * 1

	QueryRouteName  = "postgresql-query"
	InsertRouteName = "postgresql-insert"
	UpdateRouteName = "postgresql-update"
	DeleteRouteName = "postgresql-delete"
	PingRouteName   = "postgresql-ping"
)

// Attr - key value pair
type Attr struct {
	Key string
	Val any
}

// Readiness - package readiness
func Readiness() *core.Status {
	if isReady() {
		return core.StatusOK()
	}
	return core.NewStatus(core.StatusNotStarted)
}

// Query -  process a SQL select statement
func Query(ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows pgx.Rows, status *core.Status) {
	req := newQueryRequestFromValues(h, resource, template, values, args...)
	req.queryFunc = accessQuery
	return query(ctx, req)
}

// QueryT -  process a SQL select statement, returning a type
func QueryT[T Scanner[T]](ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows []T, status *core.Status) {
	req := newQueryRequestFromValues(h, resource, template, values, args...)
	req.queryFunc = accessQuery
	r, status1 := query(ctx, req)
	if !status1.OK() {
		return nil, status1
	}
	return Scan[T](r)
}

// Insert - execute a SQL insert statement
func Insert(ctx context.Context, h http.Header, resource, template string, values [][]any, args ...any) (tag CommandTag, status *core.Status) {
	req := newInsertRequest(h, resource, template, values, args...)
	return exec(ctx, req)
}

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, template string, where []Attr, args []Attr) (tag CommandTag, status *core.Status) {
	req := newUpdateRequest(h, resource, template, convert(where), convert(args))
	return exec(ctx, req)
}

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, template string, where []Attr, args ...any) (tag CommandTag, status *core.Status) {
	req := newDeleteRequest(h, resource, template, convert(where), args...)
	return exec(ctx, req)
}

// Stat - retrieve Pgx pool stats
func Stat() (*pgxpool.Stat, *core.Status) {
	return stat()
}

// Ping - ping the database cluster
func Ping(ctx context.Context, h http.Header) *core.Status {
	req := newPingRequest(h)
	return ping(ctx, req)
}
