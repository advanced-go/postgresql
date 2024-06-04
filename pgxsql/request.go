package pgxsql

import (
	"context"
	"fmt"
	"github.com/advanced-go/postgresql/pgxdml"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"time"
)

const (
	postgresScheme = "postgres"
	queryRoot      = "query"
	execRoot       = "exec"
	pingRoot       = "ping"

	selectMethod = "select"
	insertMethod = "insert"
	updateMethod = "update"
	deleteMethod = "delete"
	pingMethod   = "ping"

	selectCmd = 0
	insertCmd = 1
	updateCmd = 2
	deleteCmd = 3
	pingCmd   = 4

	nullExpectedCount = int64(-1)
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	expectedCount int64
	cmd           int
	duration      time.Duration

	resource  string
	template  string
	uri       string
	routeName string

	values    [][]any
	attrs     []pgxdml.Attr
	where     []pgxdml.Attr
	args      []any
	error     error
	header    http.Header
	queryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	execFunc  func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func newRequest(h http.Header, cmd int, resource, template, uri, routeName string, duration time.Duration) *request {
	r := new(request)
	r.expectedCount = nullExpectedCount
	r.cmd = cmd

	r.resource = resource
	r.template = template
	r.uri = uri
	r.routeName = routeName

	r.header = h

	r.duration = duration
	return r
}

func method(r *request) string {
	switch r.cmd {
	case selectCmd:
		return selectMethod
	case insertCmd:
		return insertMethod
	case updateCmd:
		return updateMethod
	case deleteCmd:
		return deleteMethod
	case pingCmd:
		return pingMethod
	}
	return "unknown"
}

func header(r *request) http.Header {
	return r.header
}

func NewHttpRequest(r *request) *http.Request {
	req, _ := http.NewRequest(method(r), r.uri, nil)
	req.Header = r.header
	return req
}

func setTimeout(ctx context.Context, req *request) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, nil
	}
	return context.WithTimeout(ctx, req.duration)
}

func buildUri(resource, path string) string {
	return fmt.Sprintf("%v://%v/%v/%v/%v", postgresScheme, "host-name", "database-name", resource, path)
	//originUrn(nid, nss, resource) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// buildQueryUri - build an uri with the Query NSS
func buildQueryUri(path string) string {
	return buildUri(queryRoot, path)
}

// buildInsertUri - build an uri with the Insert NSS
//func buildInsertUri(resource string) string {
//	return buildUri(postgresNID, insertNSS, resource)
//}

// buildUpdateUri - build an uri with the Update NSS
//func buildUpdateUri(resource string) string {
//	return buildUri(postgresNID, updateNSS, resource)
//}

// buildDeleteUri - build an uri with the Delete NSS
//func buildDeleteUri(resource string) string {
//	return buildUri(postgresNID, deleteNSS, resource)
//}

// buildFileUri - build an uri with the Query NSS
//func buildFileUri(resource string) string {
//	return buildUri(postgresNID, queryNSS, resource)
//}

func newQueryRequest(h http.Header, resource, template string, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(h, selectCmd, resource, template, buildQueryUri(resource), QueryRouteName, QueryTimeout)
	r.where = where
	r.args = args
	return r
}

func newQueryRequestFromValues(h http.Header, resource, template string, values map[string][]string, args ...any) *request {
	r := newRequest(h, selectCmd, resource, template, buildQueryUri(resource), QueryRouteName, QueryTimeout)
	r.where = buildWhere(values)
	r.args = args
	return r
}

func newInsertRequest(h http.Header, resource, template string, values [][]any, args ...any) *request {
	r := newRequest(h, insertCmd, resource, template, buildUri(execRoot, resource), InsertRouteName, InsertTimeout)
	r.values = values
	r.args = args
	return r
}

func newUpdateRequest(h http.Header, resource, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(h, updateCmd, resource, template, buildUri(execRoot, resource), UpdateRouteName, UpdateTimeout)
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
}

func newDeleteRequest(h http.Header, resource, template string, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(h, deleteCmd, resource, template, buildUri(execRoot, resource), DeleteRouteName, DeleteTimeout)
	r.where = where
	r.args = args
	return r
}

func newPingRequest(h http.Header) *request {
	r := newRequest(h, pingCmd, "", "", buildUri(pingRoot, ""), PingRouteName, PingTimeout)
	return r
}

// BuildWhere - build the []Attr based on the URL query parameters
func buildWhere(values map[string][]string) []pgxdml.Attr {
	if len(values) == 0 {
		return nil
	}
	var where []pgxdml.Attr
	for k, v := range values {
		where = append(where, pgxdml.Attr{Key: k, Val: v[0]})
	}
	return where
}

func convert(attrs []Attr) []pgxdml.Attr {
	result := make([]pgxdml.Attr, len(attrs))
	for _, pair := range attrs {
		result = append(result, pgxdml.Attr{Key: pair.Key, Val: pair.Val})
	}
	return result
}
