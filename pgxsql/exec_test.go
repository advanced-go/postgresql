package pgxsql

import (
	"errors"
	"fmt"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/postgresql/pgxdml"
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
)

func execTestProxy(req *Request) (tag CommandTag, err error) {
	switch req.Uri {
	case BuildUpdateUri(execUpdateRsc):
		return tag, errors.New("exec error")
	case BuildInsertUri(execInsertRsc):
		return CommandTag{
			Sql:          "INSERT 1",
			RowsAffected: 1234,
			Insert:       true,
			Update:       false,
			Delete:       false,
			Select:       false,
		}, nil
	}
	return tag, nil
}

func ExampleExec_Proxy() {
	ctx := runtime.ContextWithProxy(nil, execTestProxy)

	cmd, status := Exec[runtime.DebugError](ctx, NullCount, NewUpdateRequest(execUpdateRsc, execUpdateSql, nil, nil))
	fmt.Printf("test: Exec[DebugError](%v) -> %v [cmd:%v]\n", execUpdateSql, status, cmd)

	cmd, status = Exec[runtime.DebugError](ctx, NullCount, NewInsertRequest(execInsertRsc, execInsertSql, nil))
	fmt.Printf("test: Exec[DebugError](%v) -> %v [cmd:%v]\n", execInsertSql, status, cmd)

	//Output:
	//[[] github.com/gotemplates/postgresql/pgxsql/exec [exec error]]
	//test: Exec[DebugError](update test) -> Internal [cmd:{ 0 false false false false}]
	//test: Exec[DebugError](insert test) -> OK [cmd:{INSERT 1 1234 true false false false}]

}

func ExampleExec_Insert() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		cond := TestConditions{
			Time:        time.Now().UTC(),
			Location:    "plano",
			Temperature: 101.33,
		}
		req := NewInsertRequest(execInsertRsc, execInsertConditions, pgxdml.NewInsertValues([]any{pgxdml.TimestampFn, cond.Location, cond.Temperature}))

		results, status := Exec[runtime.DebugError](nil, NullCount, req)
		if !status.OK() {
			fmt.Printf("test: Insert[runtime.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execInsertConditions, status, results)
		} else {
			fmt.Printf("test: Insert[runtime.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execInsertConditions, status, results)
		}
	}

	//Output:
	//test: Insert[runtime.DebugError](nil,INSERT INTO conditions (time,location,temperature) VALUES) -> [status:OK] [cmd:{INSERT 0 1 1 true false false false}]

}

func ExampleExec_Update() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		attrs := []pgxdml.Attr{{"Temperature", 45.1234}}
		where := []pgxdml.Attr{{"Location", "plano"}}
		req := NewUpdateRequest(execUpdateRsc, execUpdateConditions, attrs, where)

		results, status := Exec[runtime.DebugError](nil, NullCount, req)
		if !status.OK() {
			fmt.Printf("test: Update[runtime.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execUpdateConditions, status, results)
		} else {
			fmt.Printf("test: Update[runtime.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execUpdateConditions, status, results)
		}
	}

	//Output:
	//test: Update[runtime.DebugError](nil,UPDATE conditions) -> [status:OK] [cmd:{UPDATE 1 1 false true false false}]

}

func ExampleExec_Delete() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		where := []pgxdml.Attr{{"Location", "plano"}}
		req := NewDeleteRequest(execDeleteRsc, execDeleteConditions, where)

		results, status := Exec[runtime.DebugError](nil, NullCount, req)
		if !status.OK() {
			fmt.Printf("test: Delete[runtime.DebugError](nil,%v) -> [status:%v] [tag:%v}\n", execDeleteConditions, status, results)
		} else {
			fmt.Printf("test: Delete[runtime.DebugError](nil,%v) -> [status:%v] [cmd:%v]\n", execDeleteConditions, status, results)
		}
	}

	//Output:
	//test: Delete[runtime.DebugError](nil,DELETE FROM conditions) -> [status:OK] [cmd:{DELETE 1 1 false false true false}]

}
