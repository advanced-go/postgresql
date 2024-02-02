package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/pgxdml"
	"net/http"
)

const (
	queryNSS  = "query"
	selectNSS = "select"
	insertNSS = "insert"
	updateNSS = "update"
	deleteNSS = "delete"
	pingNSS   = "ping"
	statNSS   = "stat"

	postgresNID = "postgresql"
	PingUri     = postgresNID + ":" + pingNSS
	StatUri     = postgresNID + ":" + statNSS

	selectCmd = 0
	insertCmd = 1
	updateCmd = 2
	deleteCmd = 3
	pingCmd   = 4

	//queryRouteName  = "query"
	//insertRouteName = "insert"
	//updateRouteName = "update"
	//deleteRouteName = "delete"

	nullExpectedCount = int64(-1)
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	expectedCount int64
	cmd           int

	resource       string
	template       string
	uri            string
	controllerName string

	values [][]any
	attrs  []pgxdml.Attr
	where  []pgxdml.Attr
	args   []any
	error  error
	header http.Header
}

func newRequest(h http.Header, cmd int, resource, template, uri, controllerName string) *request {
	r := new(request)
	r.expectedCount = nullExpectedCount
	r.cmd = cmd

	r.resource = resource
	r.template = template
	r.uri = uri
	r.controllerName = controllerName

	r.header = h
	return r
}

func method(r *request) string {
	switch r.cmd {
	case selectCmd:
		return selectNSS
	case insertCmd:
		return insertNSS
	case updateCmd:
		return updateNSS
	case deleteCmd:
		return deleteNSS
	case pingCmd:
		return pingNSS
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

func originUrn(nid, nss, resource string) string {
	return fmt.Sprintf("%v.%v.%v:%v.%v", nid, "region", "zone", nss, resource)
}

func buildUri(_, nss, resource string) string {
	return fmt.Sprintf("%v:%v.%v", PkgPath, nss, resource)
	//originUrn(nid, nss, resource) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// buildQueryUri - build an uri with the Query NSS
func buildQueryUri(resource string) string {
	return buildUri(postgresNID, queryNSS, resource)
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
	r := newRequest(h, selectCmd, resource, template, buildQueryUri(resource), queryControllerName)
	r.where = where
	r.args = args
	return r
}

func newQueryRequestFromValues(h http.Header, resource, template string, values map[string][]string, args ...any) *request {
	r := newRequest(h, selectCmd, resource, template, buildQueryUri(resource), queryControllerName)
	r.where = buildWhere(values)
	r.args = args
	return r
}

func newInsertRequest(h http.Header, resource, template string, values [][]any, args ...any) *request {
	r := newRequest(h, insertCmd, resource, template, buildUri(postgresNID, insertNSS, resource), insertControllerName)
	r.values = values
	r.args = args
	return r
}

func newUpdateRequest(h http.Header, resource, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(h, updateCmd, resource, template, buildUri(postgresNID, updateNSS, resource), updateControllerName)
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
}

func newDeleteRequest(h http.Header, resource, template string, where []pgxdml.Attr, args ...any) *request {
	r := newRequest(h, deleteCmd, resource, template, buildUri(postgresNID, deleteNSS, resource), deleteControllerName)
	r.where = where
	r.args = args
	return r
}

/*
func newPingRequest(h http.Header) *request {
	r := newRequest(h, pingCmd, pingThreshold, "", "", PingUri, pingRouteName)
	return r
}
*/

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
