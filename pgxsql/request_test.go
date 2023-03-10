package pgxsql

import (
	"fmt"
	"net/url"
)

func ExampleBuildRequest() {
	rsc := "exec-test-resource.dev"
	uri := BuildInsertUri(rsc)

	fmt.Printf("test: BuildInsertUri(%v) -> %v\n", rsc, uri)

	rsc = "query-test-resource.prod"
	uri = BuildQueryUri(rsc)

	fmt.Printf("test: BuildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: BuildInsertUri(exec-test-resource.dev) -> urn:postgres:insert.exec-test-resource.dev
	//test: BuildQueryUri(query-test-resource.prod) -> urn:postgres:query.query-test-resource.prod

}

func ExampleRequest_Validate() {
	uri := "urn:postgres:query.resource"
	sql := "select * from table"
	req := Request{}

	err := req.Validate()
	fmt.Printf("test: Validate(empty) -> %v\n", err)

	req.Uri = uri
	err = req.Validate()
	fmt.Printf("test: Validate(%v) -> %v\n", uri, err)

	req.Uri = ""
	req.Template = sql
	err = req.Validate()
	fmt.Printf("test: Validate(%v) -> %v\n", sql, err)

	req.Uri = uri
	req.Template = sql
	err = req.Validate()
	fmt.Printf("test: Validate(all) -> %v\n", err)

	//rsc := "access-log"
	//t := "delete from access_log"
	//req1 := NewDeleteRequest(rsc, t, nil)
	//err = req1.Validate()
	//fmt.Printf("test: Validate(%v) -> %v\n", t, err)

	//t = "update access_log"
	//req1 = NewUpdateRequest(rsc, t, nil, nil)
	//err = req1.Validate()
	//fmt.Printf("test: Validate(%v) -> %v\n", t, err)

	//t = "update access_log"
	//req1 = NewUpdateRequest(rsc, t, []pgxdml.Attr{{Name: "test", Val: "test"}}, nil)
	//err = req1.Validate()
	//fmt.Printf("test: Validate(%v) -> %v\n", t, err)

	//Output:
	//test: Validate(empty) -> invalid argument: request Uri is empty
	//test: Validate(urn:postgres:query.resource) -> invalid argument: request template is empty
	//test: Validate(select * from table) -> invalid argument: request Uri is empty
	//test: Validate(all) -> <nil>

}

func ExampleBuildSql() {
	rsc := "access-log"
	t := "delete from access_log"
	req := NewDeleteRequest(rsc, t, nil)

	sql := req.BuildSql()
	fmt.Printf("test: Delete.BuildSql(%v) -> %v\n", t, sql)

	t = "update access_log"
	req = NewUpdateRequest(rsc, t, nil, nil)
	sql = req.BuildSql()
	fmt.Printf("test: Update.BuildSql(%v) -> %v\n", t, sql)

	//Output:
	//test: Delete.BuildSql(delete from access_log) -> delete from access_log
	//test: Update.BuildSql(update access_log) -> update access_log

}

func ExampleNewQueryRequestFromValuesBuildSql() {
	u, _ := url.Parse("https://www.google.com/search?location=texas&zone=plano")
	rsc := "access-log"
	t := "select * from access_log {where} order by start_time desc limit 2"
	req := NewQueryRequestFromValues(rsc, t, u.Query())

	sql := req.BuildSql()
	fmt.Printf("test: NewQueryRequestFromValues(%v) -> %v\n", t, sql)

	req = NewQueryRequestFromValues(rsc, t, nil)
	sql = req.BuildSql()
	fmt.Printf("test: NewQueryRequestFromValues(%v) -> %v\n", t, sql)

	//Output:
	//test: NewQueryRequestFromValues(select * from access_log {where} order by start_time desc limit 2) -> select * from access_log
	//WHERE zone = 'plano' AND location = 'texas' order by start_time desc limit 2
	//test: NewQueryRequestFromValues(select * from access_log {where} order by start_time desc limit 2) -> select * from access_log order by start_time desc limit 2

}
