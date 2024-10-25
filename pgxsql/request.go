package pgxsql

import (
	"context"
	"fmt"
	"github.com/advanced-go/postgresql/module"
	"github.com/advanced-go/postgresql/pgxdml"
	"net/http"
	"time"
)

const (
	protocol       = "message-based"
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

	values  [][]any
	values2 map[string][]string
	attrs   []pgxdml.Attr
	where   []pgxdml.Attr
	args    []any
	error   error
	h       http.Header
}

func newRequest(cmd int, resource, template, uri, routeName string) *request {
	r := new(request)
	r.expectedCount = nullExpectedCount
	r.cmd = cmd

	r.resource = resource
	r.template = template
	r.uri = uri
	r.routeName = routeName
	r.duration = -1
	r.h = make(http.Header)
	return r
}

func (r *request) Header() http.Header {
	return r.h
}

func (r *request) Method() string {
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

func (r *request) Url() string {
	return r.uri
}

func (r *request) Protocol() string {
	return protocol
}

func (r *request) setTimeout(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if d, ok := ctx.Deadline(); ok {
		r.duration = time.Until(d)
	}
	return ctx
}

func buildUri(root, resource string) string {
	return fmt.Sprintf("%v://%v/%v:%v/%v/%v", postgresScheme, "host-name", module.Authority, "database-name", root, resource)
	//originUrn(nid, nss, test) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, test)
}

// buildQueryUri - build an uri with the Query NSS
func buildQueryUri(resource string) string {
	return buildUri(queryRoot, resource)
}

// buildInsertUri - build an uri with the Insert NSS
//func buildInsertUri(test string) string {
//	return buildUri(postgresNID, insertNSS, test)
//}

// buildUpdateUri - build an uri with the Update NSS
//func buildUpdateUri(test string) string {
//	return buildUri(postgresNID, updateNSS, test)
//}

// buildDeleteUri - build an uri with the Delete NSS
//func buildDeleteUri(test string) string {
//	return buildUri(postgresNID, deleteNSS, test)
//}

// buildFileUri - build an uri with the Query NSS
//func buildFileUri(test string) string {
//	return buildUri(postgresNID, queryNSS, test)
//}

func newQueryRequest(resource, template string, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(selectCmd, resource, template, buildQueryUri(resource), QueryRouteName)
	r.where = where
	r.args = args
	return r
}

func newQueryRequestFromValues(resource, template string, values map[string][]string, args ...any) *request {
	r := newRequest(selectCmd, resource, template, buildQueryUri(resource), QueryRouteName)
	r.where = buildWhere(values)
	r.args = args
	r.values2 = values
	return r
}

func newInsertRequest(resource, template string, values [][]any, args ...any) *request {
	r := newRequest(insertCmd, resource, template, buildUri(execRoot, resource), InsertRouteName)
	r.values = values
	r.args = args
	return r
}

func newUpdateRequest(resource, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(updateCmd, resource, template, buildUri(execRoot, resource), UpdateRouteName)
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
}

func newDeleteRequest(resource, template string, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(deleteCmd, resource, template, buildUri(execRoot, resource), DeleteRouteName)
	r.where = where
	r.args = args
	return r
}

func newPingRequest() *request {
	r := newRequest(pingCmd, "", "", buildUri(pingRoot, ""), PingRouteName)
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
