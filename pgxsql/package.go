package pgxsql

import (
	"context"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/jsonx"
	"github.com/advanced-go/postgresql/module"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	PkgPath       = "github/advanced-go/postgresql/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"

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

// QueryFunc - type declaration
type QueryFunc func(context.Context, http.Header, string, string, map[string][]string, ...any) (pgx.Rows, *core.Status)

// Query -  process a SQL select statement
func Query(ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows pgx.Rows, status *core.Status) {
	req := newQueryRequestFromValues(resource, template, values, args...)
	start := time.Now().UTC()
	rows, status = query(ctx, req)
	log(start, h, req, status)
	return rows, status
}

// QueryFuncT - type declaration
type QueryFuncT[T Scanner[T]] func(context.Context, http.Header, string, string, map[string][]string, ...any) ([]T, *core.Status)

// QueryT -  process a SQL select statement, returning a type
func QueryT[T Scanner[T]](ctx context.Context, h http.Header, resource, template string, values map[string][]string, args ...any) (rows []T, status *core.Status) {
	req := newQueryRequestFromValues(resource, template, values, args...)
	req.Header().Set(core.XTo, module.Authority)
	start := time.Now().UTC()
	_, resp, status1 := core.ExchangeHeaders(h)
	if resp != "" || status1 != "" {
		status = core.StatusOK()
		ctx = req.setTimeout(ctx)
		if resp != "" {
			rows, status = Unmarshal[T](resp)
		}
		if status1 != "" {
			status = jsonx.NewStatusFrom(status1)
		}
		log(start, h, req, status)
		return
	}
	r, status2 := query(ctx, req)
	log(start, h, req, status2)
	if !status2.OK() {
		return nil, status2
	}
	return Scan[T](r)
}

// InsertFunc - type
type InsertFunc func(context.Context, http.Header, string, string, [][]any, ...any) (CommandTag, *core.Status)

// Insert - execute a SQL insert statement
func Insert(ctx context.Context, h http.Header, resource, template string, values [][]any, args ...any) (tag CommandTag, status *core.Status) {
	req := newInsertRequest(resource, template, values, args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// InsertFuncT - type
type InsertFuncT[T Scanner[T]] func(context.Context, http.Header, string, string, []T, ...any) (CommandTag, *core.Status)

// InsertT - execute a SQL insert statement
func InsertT[T Scanner[T]](ctx context.Context, h http.Header, resource, template string, entries []T, args ...any) (tag CommandTag, status *core.Status) {
	_, _, stat1 := core.ExchangeHeaders(h)
	if stat1 != "" {
		start := time.Now().UTC()
		req := newInsertRequest(resource, template, nil, args...)
		ctx = req.setTimeout(ctx)
		status = jsonx.NewStatusFrom(stat1)
		log(start, h, req, status)
		return
	}
	rows, status1 := Rows[T](entries)
	if !status1.OK() {
		return CommandTag{}, status1
	}
	req := newInsertRequest(resource, template, rows, args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// UpdateFunc - type
type UpdateFunc func(context.Context, http.Header, string, string, []Attr, []Attr) (CommandTag, *core.Status)

// Update - execute a SQL update statement
func Update(ctx context.Context, h http.Header, resource, template string, where []Attr, args []Attr) (tag CommandTag, status *core.Status) {
	req := newUpdateRequest(resource, template, convert(where), convert(args))
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// DeleteFunc - type
type DeleteFunc func(context.Context, http.Header, string, string, []Attr, ...any) (CommandTag, *core.Status)

// Delete - execute a SQL delete statement
func Delete(ctx context.Context, h http.Header, resource, template string, where []Attr, args ...any) (tag CommandTag, status *core.Status) {
	req := newDeleteRequest(resource, template, convert(where), args...)
	start := time.Now().UTC()
	tag, status = exec(ctx, req)
	log(start, h, req, status)
	return tag, status
}

// Stat - retrieve Pgx pool stats
func Stat() (*pgxpool.Stat, *core.Status) {
	return stat()
}

// Ping - ping the database cluster
func Ping(ctx context.Context, h http.Header) *core.Status {
	req := newPingRequest()
	start := time.Now().UTC()
	status := ping(ctx, req)
	log(start, h, req, status)
	return status
}
