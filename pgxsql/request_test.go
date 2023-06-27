package pgxsql

import (
	"fmt"
)

func ExampleBuildRequest() {
	rsc := "exec-test-resource.dev"
	uri := BuildInsertUri(rsc)

	fmt.Printf("test: BuildInsertUri(%v) -> %v\n", rsc, uri)

	rsc = "query-test-resource.prod"
	uri = BuildQueryUri(rsc)

	fmt.Printf("test: BuildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: BuildInsertUri(exec-test-resource.dev) -> urn:postgresql.region.zone:insert.exec-test-resource.dev
	//test: BuildQueryUri(query-test-resource.prod) -> urn:postgresql.region.zone:query.query-test-resource.prod

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
