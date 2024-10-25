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
	rsc := "test-test.dev"
	uri := buildUri(execRoot, rsc)
	fmt.Printf("test: buildInsertUri(%v) -> %v\n", rsc, uri)

	rsc = "test-test.prod"
	uri = buildQueryUri(rsc)
	fmt.Printf("test: buildQueryUri(%v) -> %v\n", rsc, uri)

	//Output:
	//test: buildInsertUri(test-test.dev) -> postgres://host-name/github/advanced-go/postgresql:database-name/exec/test-test.dev
	//test: buildQueryUri(test-test.prod) -> postgres://host-name/github/advanced-go/postgresql:database-name/query/test-test.prod

}

func ExampleRequest_Validate() {
	uri := "urn:postgres:query.test"
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
	//test: Validate(urn:postgres:query.test) -> invalid argument: request template is empty
	//test: Validate(select * from table) -> invalid argument: request Uri is empty
	//test: Validate(all) -> <nil>

}
