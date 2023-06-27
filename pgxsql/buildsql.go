package pgxsql

import (
	"github.com/go-ai-agent/postgresql/pgxdml"
)

func BuildSql(r *Request) string {
	var stmt = r.Template
	var err error

	switch r.Cmd {
	case SelectCmd:
		stmt, err = pgxdml.ExpandSelect(r.Template, r.Where)
	case InsertCmd:
		if len(r.Values) > 0 {
			stmt, err = pgxdml.WriteInsert(r.Template, r.Values)
		}
	case UpdateCmd:
		//if len(r.Where) == 0 {
		//	r.Where = append(r.Where, pgxdml.Attr{Name: "update_error_no_where_clause", Val: "null"})
		//}
		//if len(r.Attrs) == 0 {
		//	r.Attrs = append(r.Attrs, pgxdml.Attr{Name: "update_error_no_set_clause", Val: "null"})
		//}
		if len(r.Where) > 0 && len(r.Attrs) > 0 {
			stmt, err = pgxdml.WriteUpdate(r.Template, r.Attrs, r.Where)
		}
	case DeleteCmd:
		if len(r.Where) > 0 {
			//r.Where = append(r.Where, pgxdml.Attr{Name: "delete_error_no_where_clause", Val: "null"})
			stmt, err = pgxdml.WriteDelete(r.Template, r.Where)
		}
	}
	r.Error = err
	return stmt
}
