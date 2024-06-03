package pgxsql

import (
	"context"
	"github.com/advanced-go/stdlib/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type pkg struct{}

const (
	PkgPath       = "github/advanced-go/postgresql/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"
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
	return query(ctx, req)
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

// Stat - retrieve core stats
func Stat(ctx context.Context) (*pgxpool.Stat, *core.Status) {
	return stat(ctx)
}

// Ping - ping the database cluster
func Ping(ctx context.Context) *core.Status {
	return ping(ctx)
}
