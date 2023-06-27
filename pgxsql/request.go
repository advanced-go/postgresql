package pgxsql

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

const (
	QueryNSS  = "query"
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

func BuildUri(nid string, o runtime.Origin, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return BuildUri(PostgresNID, origin, QueryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return BuildUri(PostgresNID, origin, InsertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return BuildUri(PostgresNID, origin, UpdateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return BuildUri(PostgresNID, origin, DeleteNSS, resource)
}

// Request - contains data needed to build the SQL statement related to the uri
type Request struct {
	ExpectedCount int64
	Cmd           int
	Uri           string
	Template      string
	Values        [][]any
	Attrs         []runtime.Attr
	Where         []runtime.Attr
	Args          []any
	Error         error
}

func (r *Request) Validate() error {
	if r.Uri == "" {
		return errors.New("invalid argument: request Uri is empty")
	}
	if r.Template == "" {
		return errors.New("invalid argument: request template is empty")
	}
	return nil
}

func (r *Request) String() string {
	return r.Template
}

func NewQueryRequest(uri, template string, where []runtime.Attr, args ...any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: where, Args: args}
}

func NewQueryRequestFromValues(uri, template string, values map[string][]string, args ...any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: BuildWhere(values), Args: args}
}

func NewInsertRequest(uri, template string, values [][]any, args ...any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: InsertCmd, Uri: uri, Template: template, Values: values, Args: args}
}

func NewUpdateRequest(uri, template string, attrs []runtime.Attr, where []runtime.Attr, args ...any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: UpdateCmd, Uri: uri, Template: template, Attrs: attrs, Where: where, Args: args}
}

func NewDeleteRequest(uri, template string, where []runtime.Attr, args ...any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: DeleteCmd, Uri: uri, Template: template, Attrs: nil, Where: where, Args: args}
}

// BuildWhere - build the []Attr based on the URL query parameters
func BuildWhere(values map[string][]string) []runtime.Attr {
	if len(values) == 0 {
		return nil
	}
	var where []runtime.Attr
	for k, v := range values {
		where = append(where, runtime.Attr{Key: k, Val: v[0]})
	}
	return where
}
