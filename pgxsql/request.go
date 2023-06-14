package pgxsql

import (
	"fmt"
	"github.com/go-sre/core/sql"
	"github.com/go-sre/postgresql/pgxdml"
)

const (
	PostgresNID = "postgresql"
	//QueryNSS    = "query"
	//InsertNSS   = "insert"
	//UpdateNSS   = "update"
	//DeleteNSS   = "delete"
	//PingNSS     = "ping"
	//StatNSS     = "stat"

	PingUri = "urn:" + PostgresNID + ":" + sql.PingNSS
	StatUri = "urn:" + PostgresNID + ":" + sql.StatNSS

	//selectCmd = 0
	//insertCmd = 1
	//updateCmd = 2
	//deleteCmd = 3

	variableReference = "$1"
)

func buildUri(nsid, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nsid, region, zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return sql.BuildQueryUri(PostgresNID, region, zone, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return sql.BuildInsertUri(PostgresNID, region, zone, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return sql.BuildUpdateUri(PostgresNID, region, zone, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return sql.BuildDeleteUri(PostgresNID, region, zone, resource)
}

// Request - contains data needed to build the SQL statement related to the uri
/*
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

	return nil
}

func (r *Request) String() string {
	return r.Template
}


*/

func BuildSql(r *sql.Request) string {
	var stmt = r.Template
	var err error

	switch r.Cmd {
	case sql.SelectCmd:
		stmt, err = pgxdml.ExpandSelect(r.Template, r.Where)
	case sql.InsertCmd:
		if len(r.Values) > 0 {
			stmt, err = pgxdml.WriteInsert(r.Template, r.Values)
		}
	case sql.UpdateCmd:
		//if len(r.Where) == 0 {
		//	r.Where = append(r.Where, pgxdml.Attr{Name: "update_error_no_where_clause", Val: "null"})
		//}
		//if len(r.Attrs) == 0 {
		//	r.Attrs = append(r.Attrs, pgxdml.Attr{Name: "update_error_no_set_clause", Val: "null"})
		//}
		if len(r.Where) > 0 && len(r.Attrs) > 0 {
			stmt, err = pgxdml.WriteUpdate(r.Template, r.Attrs, r.Where)
		}
	case sql.DeleteCmd:
		if len(r.Where) > 0 {
			//r.Where = append(r.Where, pgxdml.Attr{Name: "delete_error_no_where_clause", Val: "null"})
			stmt, err = pgxdml.WriteDelete(r.Template, r.Where)
		}
	}
	r.Error = err
	return stmt
}

func NewQueryRequest(resource, template string, where []sql.Attr) *sql.Request {
	return sql.NewQueryRequest(BuildQueryUri(resource), template, where)
}

func NewQueryRequestFromValues(resource, template string, values map[string][]string) *sql.Request {
	return sql.NewQueryRequestFromValues(BuildQueryUri(resource), template, values)
}

func NewInsertRequest(resource, template string, values [][]any) *sql.Request {
	return sql.NewInsertRequest(BuildInsertUri(resource), template, values)
}

func NewUpdateRequest(resource, template string, attrs []sql.Attr, where []sql.Attr) *sql.Request {
	return sql.NewUpdateRequest(BuildUpdateUri(resource), template, attrs, where)
}

func NewDeleteRequest(resource, template string, where []sql.Attr) *sql.Request {
	return sql.NewDeleteRequest(BuildDeleteUri(resource), template, where)
}
