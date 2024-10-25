package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/pgxdml"
	"time"
)

func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

const (
	execUpdateSql = "update test"
	execInsertSql = "insert test"
	execUpdateRsc = "update"
	execInsertRsc = "insert"
	execDeleteRsc = "delete"

	execInsertConditions = "INSERT INTO conditions (time,location,temperature) VALUES"
	execUpdateConditions = "UPDATE conditions"
	execDeleteConditions = "DELETE FROM conditions"

	status504    = "file://[cwd]/test/status-504.json"
	updateCmdTag = "file://[cwd]/test/update-cmd-tag.json"
)

func ExampleExec_Status() {
	//lookup.SetOverride(status504)
	result, status := exec(nil, newUpdateRequest(execUpdateRsc, execUpdateSql, nil, nil))
	fmt.Printf("test: Exec(ctx,%v) -> [tag:%v] [status:%v]\n", execUpdateSql, result, status)

	//Output:
	//test: Exec(ctx,update test) -> [tag:{ 0 false false false false}] [status:Timeout [error 1]]

}

func ExampleExec_Proxy() {
	//lookup.SetOverride(updateCmdTag)
	req := newUpdateRequest(execUpdateRsc, execUpdateSql, nil, nil)
	tag, status := exec(nil, req)
	fmt.Printf("test: Exec(%v) -> [cmd:%v] [status:%v]\n", execUpdateSql, tag, status)

	//Output:
	//test: Exec(update test) -> [cmd:{update test 0 false true false false}] [status:OK]

}

func ExampleExec_Insert() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		cond := TestConditions{
			Time:        time.Now().UTC(),
			Location:    "plano",
			Temperature: 101.33,
		}
		req := newInsertRequest(execInsertRsc, execInsertConditions, pgxdml.NewInsertValues([]any{pgxdml.TimestampFn, cond.Location, cond.Temperature}))

		results, status := exec(nil, req)
		if !status.OK() {
			fmt.Printf("test: Insert(nil,%v) -> [status:%v] [tag:%v}\n", execInsertConditions, status, results)
		} else {
			fmt.Printf("test: Insert(nil,%v) -> [status:%v] [cmd:%v]\n", execInsertConditions, status, results)
		}
	}

	//Output:
	//test: Insert(nil,INSERT INTO conditions (time,location,temperature) VALUES) -> [status:OK] [cmd:{INSERT 0 1 1 true false false false}]

}

func ExampleExec_Update() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		attrs := []pgxdml.Attr{{"Temperature", 45.1234}}
		where := []pgxdml.Attr{{"Location", "plano"}}
		req := newUpdateRequest(execUpdateRsc, execUpdateConditions, attrs, where)

		results, status := exec(nil, req)
		if !status.OK() {
			fmt.Printf("test: Update(nil,%v) -> [status:%v] [tag:%v}\n", execUpdateConditions, status, results)
		} else {
			fmt.Printf("test: Update(nil,%v) -> [status:%v] [cmd:%v]\n", execUpdateConditions, status, results)
		}
	}

	//Output:
	//test: Update(nil,UPDATE conditions) -> [status:OK] [cmd:{UPDATE 1 1 false true false false}]

}

func ExampleExec_Delete() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		where := []pgxdml.Attr{{"Location", "plano"}}
		req := newDeleteRequest(execDeleteRsc, execDeleteConditions, where)

		results, status := exec(nil, req)
		if !status.OK() {
			fmt.Printf("test: Delete(nil,%v) -> [status:%v] [tag:%v}\n", execDeleteConditions, status, results)
		} else {
			fmt.Printf("test: Delete(nil,%v) -> [status:%v] [cmd:%v]\n", execDeleteConditions, status, results)
		}
	}

	//Output:
	//test: Delete(nil,DELETE FROM conditions) -> [status:OK] [cmd:{DELETE 1 1 false false true false}]

}
