package pgxsql

import (
	"errors"
	"fmt"
)

func validate(r *request) error {
	if r.uri == "" {
		return errors.New("invalid argument: request Uri is empty")
	}
	if r.template == "" {
		return errors.New("invalid argument: request template is empty")
	}
	return nil
}

func ExampleBuildRequest() {
	rsc := "exec-test-resource.dev"
	uri := buildInsertUri(rsc)

	fmt.Printf("test: buildInsertUri(%v) -> %v\n", rsc, uri)

	rsc = "query-test-resource.prod"
	uri = buildQueryUri(rsc)

	fmt.Printf("test: buildQueryUri(%v) -> %v\n", rsc, uri)

	rsc = "file://[cwd]/example-domain/activitytest/test.json"
	uri = buildQueryUri(rsc)

	fmt.Printf("test: buildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: buildInsertUri(exec-test-resource.dev) -> postgresql.region.zone:insert.exec-test-resource.dev
	//test: buildQueryUri(query-test-resource.prod) -> postgresql.region.zone:query.query-test-resource.prod
	//test: buildQueryUri(file://[cwd]/example-domain/activitytest/test.json) -> postgresql.region.zone:query.file://[cwd]/example-domain/activitytest/test.json

}

func ExampleRequest_Validate() {
	uri := "urn:postgres:query.resource"
	sql := "select * from table"
	req := request{}

	err := validate(&req)
	fmt.Printf("test: Validate(empty) -> %v\n", err)

	req.uri = uri
	err = validate(&req)
	fmt.Printf("test: Validate(%v) -> %v\n", uri, err)

	req.uri = ""
	req.template = sql
	err = validate(&req)
	fmt.Printf("test: Validate(%v) -> %v\n", sql, err)

	req.uri = uri
	req.template = sql
	err = validate(&req)
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

func ExampleNewRequest_File() {
	r := NewQueryRequest(nil, "file://[cwd]/example-domain/activitytest/test.json", "", nil)
	fmt.Printf("test: NewQueryRequest() -> %v\n", r.Uri())

	r = NewInsertRequest(nil, "file://[cwd]/example-domain/activitytest/test.json", "", nil)
	fmt.Printf("test: NewInsertRequest() -> %v\n", r.Uri())

	r = NewUpdateRequest(nil, "file://[cwd]/example-domain/activitytest/test.json", "", nil, nil)
	fmt.Printf("test: NewUpdateRequest() -> %v\n", r.Uri())

	r = NewDeleteRequest(nil, "file://[cwd]/example-domain/activitytest/test.json", "", nil, nil)
	fmt.Printf("test: NewDeleteRequest() -> %v\n", r.Uri())

	//Output:
	//test: NewQueryRequest() -> file://[cwd]/example-domain/activitytest/test.json
	//test: NewInsertRequest() -> file://[cwd]/example-domain/activitytest/test.json
	//test: NewUpdateRequest() -> file://[cwd]/example-domain/activitytest/test.json
	//test: NewDeleteRequest() -> file://[cwd]/example-domain/activitytest/test.json

}
