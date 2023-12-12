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
	PkgPath       = "github.com/advanced-go/postgresql/pgxsql"
	ReadinessPath = PkgPath + ":Readiness"
)

// Readiness - package readiness
func Readiness(uri string) runtime.Status {
	if uri == ReadinessPath {
		if isReady() {
			return runtime.StatusOK()
		}
		return runtime.NewStatus(runtime.StatusNotStarted)
	}
	return runtime.NewStatus(http.StatusNotFound)
}

// Query - function for a query
func Query(ctx context.Context, req Request) (rows pgx.Rows, status runtime.Status) {
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(req.Header(), req.Method(), req.Uri()), "Query", -1, "", access.NewStatusCodeClosure(&status))()
	return query(ctx, req)
}

// Exec - function for executing a SQL statement
func Exec(ctx context.Context, req Request) (tag CommandTag, status runtime.Status) {
	defer access.LogDeferred(access.InternalTraffic, access.NewRequest(req.Header(), req.Method(), req.Uri()), "Exec", -1, "", access.NewStatusCodeClosure(&status))()
	return exec(ctx, req)
}

// Stat - function for retrieving runtime stats
func Stat(ctx context.Context) (*pgxpool.Stat, runtime.Status) {
	return stat(ctx)
}

// Ping - function for pinging the database cluster
func Ping(ctx context.Context) runtime.Status {
	return ping(ctx)
}
