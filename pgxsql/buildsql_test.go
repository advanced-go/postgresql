package pgxsql

import "fmt"

func _ExampleBuildSql() {
	rsc := "access-log"
	t := "delete from access_log"
	req := NewDeleteRequest(nil, rsc, t, nil)

	sql := req.Sql() //BuildSql(req)
	fmt.Printf("test: Delete.BuildSql(%v) -> %v\n", t, sql)

	t = "update access_log"
	req = NewUpdateRequest(nil, rsc, t, nil, nil)
	sql = req.Sql() //BuildSql(req)
	fmt.Printf("test: Update.BuildSql(%v) -> %v\n", t, sql)

	//Output:
	//test: Delete.BuildSql(delete from access_log) -> delete from access_log
	//WHERE delete_error_no_where_clause = 'null';
	//test: Update.BuildSql(update access_log) -> update access_log
	//SET update_error_no_set_clause = 'null'
	//WHERE update_error_no_where_clause = 'null';

}
