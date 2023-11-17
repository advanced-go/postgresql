package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/pgxdml"
	"net/http"
)

const (
	QueryNSS  = "query"
	SelectNSS = "select"
	InsertNSS = "insert"
	UpdateNSS = "update"
	DeleteNSS = "delete"
	PingNSS   = "ping"
	StatNSS   = "stat"

	PostgresNID = "postgresql"
	PingUri     = "urn:" + PostgresNID + ":" + PingNSS
	StatUri     = "urn:" + PostgresNID + ":" + StatNSS

	SelectCmd = 0
	InsertCmd = 1
	UpdateCmd = 2
	DeleteCmd = 3

	NullExpectedCount = int64(-1)
)

type Request interface {
	Uri() string
	Method() string
	Sql() string
	Args() []any
	String() string
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
}

func (r *request) Uri() string {
	return r.uri
}

func (r *request) Method() string {
	switch r.cmd {
	case SelectCmd:
		return SelectNSS
	case InsertCmd:
		return InsertNSS
	case UpdateCmd:
		return UpdateNSS
	case DeleteCmd:
		return DeleteNSS
	}
	return "unknown"
}

func (r *request) Sql() string {
	return BuildSql(r)
}

func (r *request) Args() []any {
	return r.args
}

func (r *request) String() string {
	return r.template
}

func (r *request) HttpRequest() *http.Request {
	req, _ := http.NewRequest(r.Method(), r.Uri(), nil)
	return req
}

func OriginUrn(nid, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, "region", "zone", nss, resource)
}

func BuildUri(nid string, nss, resource string) string {
	return OriginUrn(nid, nss, resource) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return BuildUri(PostgresNID, QueryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return BuildUri(PostgresNID, InsertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return BuildUri(PostgresNID, UpdateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return BuildUri(PostgresNID, DeleteNSS, resource)
}

func NewQueryRequest(uri, template string, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.expectedCount = NullExpectedCount
	r.cmd = SelectCmd
	r.uri = uri
	r.template = template
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount:NullExpectedCount , Cmd: SelectCmd, Uri: uri, Template: template, Where: where, Args: args}
}

func NewQueryRequestFromValues(uri, template string, values map[string][]string, args ...any) Request {
	r := new(request)
	r.expectedCount = NullExpectedCount
	r.cmd = SelectCmd
	r.uri = uri
	r.template = template
	r.where = BuildWhere(values)
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: BuildWhere(values), Args: args}
}

func NewInsertRequest(uri, template string, values [][]any, args ...any) Request {
	r := new(request)
	r.expectedCount = NullExpectedCount
	r.cmd = InsertCmd
	r.uri = uri
	r.template = template
	r.values = values
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: InsertCmd, Uri: uri, Template: template, Values: values, Args: args}
}

func NewUpdateRequest(uri, template string, attrs []pgxdml.Attr, where []pgxdml.Attr, args ...any) Request {
	r := new(request)
	r.expectedCount = NullExpectedCount
	r.cmd = UpdateCmd
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
	r.expectedCount = NullExpectedCount
	r.cmd = DeleteCmd
	r.uri = uri
	r.template = template
	r.attrs = nil
	r.where = where
	r.args = args
	return r
	//return &Request{ExpectedCount: NullExpectedCount, Cmd: DeleteCmd, Uri: uri, Template: template, Attrs: nil, Where: where, Args: args}
}

// BuildWhere - build the []Attr based on the URL query parameters
func BuildWhere(values map[string][]string) []pgxdml.Attr {
	if len(values) == 0 {
		return nil
	}
	var where []pgxdml.Attr
	for k, v := range values {
		where = append(where, pgxdml.Attr{Key: k, Val: v[0]})
	}
	return where
}
