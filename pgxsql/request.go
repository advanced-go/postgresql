package pgxsql

import (
	"errors"
	"fmt"
	"github.com/go-sre/core/sql"
	"github.com/go-sre/postgresql/pgxdml"
)

const (
	PostgresNID = "postgresql"
	QueryNSS    = "query"
	InsertNSS   = "insert"
	UpdateNSS   = "update"
	DeleteNSS   = "delete"
	PingNSS     = "ping"
	StatNSS     = "stat"

	PingUri = "urn:" + PostgresNID + ":" + PingNSS
	StatUri = "urn:" + PostgresNID + ":" + StatNSS

	selectCmd = 0
	insertCmd = 1
	updateCmd = 2
	deleteCmd = 3

	variableReference = "$1"
)

func buildUri(nsid, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nsid, region, zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return buildUri(PostgresNID, QueryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return buildUri(PostgresNID, InsertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return buildUri(PostgresNID, UpdateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return buildUri(PostgresNID, DeleteNSS, resource)
}

// Request - contains data needed to build the SQL statement related to the uri
type Request struct {
	cmd      int
	Uri      string
	Template string
	Values   [][]any
	Attrs    []sql.Attr
	Where    []sql.Attr
	Error    error
}

func (r *Request) Validate() error {
	if r.Uri == "" {
		return errors.New("invalid argument: request Uri is empty")
	}
	if r.Template == "" {
		return errors.New("invalid argument: request template is empty")
	}
	/*if r.cmd == deleteCmd {
		if len(r.Where) == 0 {
			return errors.New("invalid argument: delete where clause is empty")
		}
	}
	if r.cmd == updateCmd {
		if len(r.Attrs) == 0 {
			return errors.New("invalid argument: update set clause is empty")
		}
		if len(r.Where) == 0 {
			return errors.New("invalid argument: update where clause is empty")
		}
	}
	*/
	return nil
}

func (r *Request) String() string {
	return r.Template
}

func (r *Request) BuildSql() string {
	var sql = r.Template
	var err error

	switch r.cmd {
	case selectCmd:
		sql, err = pgxdml.ExpandSelect(r.Template, r.Where)
	case insertCmd:
		if len(r.Values) > 0 {
			sql, err = pgxdml.WriteInsert(r.Template, r.Values)
		}
	case updateCmd:
		//if len(r.Where) == 0 {
		//	r.Where = append(r.Where, pgxdml.Attr{Name: "update_error_no_where_clause", Val: "null"})
		//}
		//if len(r.Attrs) == 0 {
		//	r.Attrs = append(r.Attrs, pgxdml.Attr{Name: "update_error_no_set_clause", Val: "null"})
		//}
		if len(r.Where) > 0 && len(r.Attrs) > 0 {
			sql, err = pgxdml.WriteUpdate(r.Template, r.Attrs, r.Where)
		}
	case deleteCmd:
		if len(r.Where) > 0 {
			//r.Where = append(r.Where, pgxdml.Attr{Name: "delete_error_no_where_clause", Val: "null"})
			sql, err = pgxdml.WriteDelete(r.Template, r.Where)
		}
	}
	r.Error = err
	return sql
}

func NewQueryRequest(resource, template string, where []sql.Attr) *Request {
	return &Request{cmd: selectCmd, Uri: BuildQueryUri(resource), Template: template, Where: where}
}

func NewQueryRequestFromValues(resource, template string, values map[string][]string) *Request {
	return &Request{cmd: selectCmd, Uri: BuildQueryUri(resource), Template: template, Where: pgxdml.BuildWhere(values)}
}

func NewInsertRequest(resource, template string, values [][]any) *Request {
	return &Request{cmd: insertCmd, Uri: BuildInsertUri(resource), Template: template, Values: values}
}

func NewUpdateRequest(resource, template string, attrs []sql.Attr, where []sql.Attr) *Request {
	return &Request{cmd: updateCmd, Uri: BuildUpdateUri(resource), Template: template, Attrs: attrs, Where: where}
}

func NewDeleteRequest(resource, template string, where []sql.Attr) *Request {
	return &Request{cmd: deleteCmd, Uri: BuildDeleteUri(resource), Template: template, Attrs: nil, Where: where}
}
