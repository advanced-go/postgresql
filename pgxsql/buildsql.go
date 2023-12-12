package pgxsql

import (
	"github.com/advanced-go/postgresql/pgxdml"
)

func buildSql(r *request) string {
	var stmt = r.template
	var err error

	switch r.cmd {
	case selectCmd:
		stmt, err = pgxdml.ExpandSelect(r.template, r.where)
	case insertCmd:
		if len(r.values) > 0 {
			stmt, err = pgxdml.WriteInsert(r.template, r.values)
		}
	case updateCmd:
		//if len(r.Where) == 0 {
		//	r.Where = append(r.Where, pgxdml.Attr{Name: "update_error_no_where_clause", Val: "null"})
		//}
		//if len(r.Attrs) == 0 {
		//	r.Attrs = append(r.Attrs, pgxdml.Attr{Name: "update_error_no_set_clause", Val: "null"})
		//}
		if len(r.where) > 0 && len(r.attrs) > 0 {
			stmt, err = pgxdml.WriteUpdate(r.template, r.attrs, r.where)
		}
	case deleteCmd:
		if len(r.where) > 0 {
			//r.Where = append(r.Where, pgxdml.Attr{Name: "delete_error_no_where_clause", Val: "null"})
			stmt, err = pgxdml.WriteDelete(r.template, r.where)
		}
	}
	r.error = err
	return stmt
}
