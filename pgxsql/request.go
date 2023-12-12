package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/pgxdml"
	"net/http"
	"strings"
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
	PingUri     = "urn:" + postgresNID + ":" + pingNSS
	StatUri     = "urn:" + postgresNID + ":" + statNSS

	fileScheme = "file://"
	selectCmd  = 0
	insertCmd  = 1
	updateCmd  = 2
	deleteCmd  = 3

	NullExpectedCount = int64(-1)
)

type Request interface {
	Uri() string
	IsFileScheme() bool
	Method() string
	Sql() string
	Args() []any
	String() string
	Header() http.Header
	HttpRequest() *http.Request
}

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	expectedCount int64
	cmd           int
	uri           string
	template      string
	values        [][]any
	attrs         []pgxdml.Attr
	where         []pgxdml.Attr
	args          []any
	error         error
	header        http.Header
	//exec          func(Request) (CommandTag, runtime.Status)
	//query         func(Request) (pgx.Rows, runtime.Status)
}

func (r *request) Uri() string {
	return r.uri
}

func (r *request) IsFileScheme() bool {
	return strings.HasPrefix(r.uri, fileScheme)
}

func (r *request) Method() string {
	switch r.cmd {
	case selectCmd:
		return selectNSS
	case insertCmd:
		return insertNSS
	case updateCmd:
		return updateNSS
	case deleteCmd:
		return deleteNSS
	}
	return "unknown"
}

func (r *request) Sql() string {
	return buildSql(r)
}

func (r *request) Args() []any {
	return r.args
}

func (r *request) String() string {
	return r.template
}

func (r *request) Header() http.Header {
	return r.header
}

func (r *request) HttpRequest() *http.Request {
	req, _ := http.NewRequest(r.Method(), r.Uri(), nil)
	req.Header = r.header
	return req
}

/*
	func (r *request) setExecProxy(proxy func(Request) (CommandTag, runtime.Status)) {
		r.exec = proxy
	}

	func (r *request) execProxy() func(Request) (CommandTag, runtime.Status) {
		return r.exec
	}

	func (r *request) setQueryProxy(proxy func(Request) (pgx.Rows, runtime.Status)) {
		r.query = proxy
	}

	func (r *request) queryProxy() func(Request) (pgx.Rows, runtime.Status) {
		return r.query
	}
*/
func originUrn(nid, nss, resource string) string {
	return fmt.Sprintf("%v.%v.%v:%v.%v", nid, "region", "zone", nss, resource)
}
func originFileUrn(nid, resource string) string {
	return fmt.Sprintf("%v.%v.%v:%v", nid, "region", "zone", resource)
}

func buildUri(nid string, nss, resource string) string {
	return originUrn(nid, nss, resource) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// buildQueryUri - build an uri with the Query NSS
func buildQueryUri(resource string) string {
	return buildUri(postgresNID, queryNSS, resource)
}

// buildInsertUri - build an uri with the Insert NSS
func buildInsertUri(resource string) string {
	return buildUri(postgresNID, insertNSS, resource)
}

// buildUpdateUri - build an uri with the Update NSS
func buildUpdateUri(resource string) string {
	return buildUri(postgresNID, updateNSS, resource)
}

// buildDeleteUri - build an uri with the Delete NSS
func buildDeleteUri(resource string) string {
	return buildUri(postgresNID, deleteNSS, resource)
}

// buildFileUri - build an uri with the Query NSS
func buildFileUri(resource string) string {
	return buildUri(postgresNID, queryNSS, resource)
}

func NewQueryRequest(resource, template string, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = selectCmd
	if strings.HasPrefix(resource, fileScheme) {
		r.uri = resource
	} else {
		r.uri = buildQueryUri(resource)
	}
	r.template = template
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount:NullExpectedCount , Cmd: SelectCmd, Uri: uri, Template: template, Where: where, Args: args}
}

func NewQueryRequestFromValues(resource, template string, values map[string][]string, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = selectCmd
	if strings.HasPrefix(resource, fileScheme) {
		r.uri = resource
	} else {
		r.uri = buildQueryUri(resource)
	}
	r.template = template
	r.where = buildWhere(values)
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: BuildWhere(values), Args: args}
}

func NewInsertRequest(resource, template string, values [][]any, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = insertCmd
	if strings.HasPrefix(resource, fileScheme) {
		r.uri = resource
	} else {
		r.uri = buildInsertUri(resource)
	}
	r.template = template
	r.values = values
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: InsertCmd, Uri: uri, Template: template, Values: values, Args: args}
}

func NewUpdateRequest(resource, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = updateCmd
	if strings.HasPrefix(resource, fileScheme) {
		r.uri = resource
	} else {
		r.uri = buildUpdateUri(resource)
	}
	r.template = template
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: UpdateCmd, Uri: uri, Template: template, Attrs: attrs, Where: where, Args: args}
}

func NewDeleteRequest(resource, template string, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = deleteCmd
	if strings.HasPrefix(resource, fileScheme) {
		r.uri = resource
	} else {
		r.uri = buildDeleteUri(resource)
	}
	r.template = template
	r.attrs = nil
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: DeleteCmd, Uri: uri, Template: template, Attrs: nil, Where: where, Args: args}
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
