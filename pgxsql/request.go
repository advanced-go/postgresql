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
	PingUri     = "urn:" + postgresNID + ":" + pingNSS
	StatUri     = "urn:" + postgresNID + ":" + statNSS

	selectCmd = 0
	insertCmd = 1
	updateCmd = 2
	deleteCmd = 3

	NullExpectedCount = int64(-1)
)

type Request interface {
	Uri() string
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
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, "region", "zone", nss, resource)
}

func buildUri(nid string, nss, resource string) string {
	return originUrn(nid, nss, resource) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return buildUri(postgresNID, queryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return buildUri(postgresNID, insertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return buildUri(postgresNID, updateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return buildUri(postgresNID, deleteNSS, resource)
}

func NewQueryRequest(uri, template string, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = selectCmd
	r.uri = uri
	r.template = template
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount:NullExpectedCount , Cmd: SelectCmd, Uri: uri, Template: template, Where: where, Args: args}
}

func NewQueryRequestFromValues(uri, template string, values map[string][]string, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = selectCmd
	r.uri = uri
	r.template = template
	r.where = buildWhere(values)
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: BuildWhere(values), Args: args}
}

func NewInsertRequest(uri, template string, values [][]any, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = insertCmd
	r.uri = uri
	r.template = template
	r.values = values
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: InsertCmd, Uri: uri, Template: template, Values: values, Args: args}
}

func NewUpdateRequest(uri, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = updateCmd
	r.uri = uri
	r.template = template
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: UpdateCmd, Uri: uri, Template: template, Attrs: attrs, Where: where, Args: args}
}

func NewDeleteRequest(uri, template string, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.header = make(http.Header)
	r.expectedCount = NullExpectedCount
	r.cmd = deleteCmd
	r.uri = uri
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
