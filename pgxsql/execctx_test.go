package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/gotemplates/core/runtime"
)

const (
	execTestUpdateSql = "update test"
	execTestInsertSql = "insert test"
	execTestUpdateRsc = "update"
	execTestInsertRsc = "insert"
)

func isExecContext(ctx context.Context) bool {
	if _, ok := any(ctx).(execWithValue); ok {
		return true
	}
	return false
}

var execCtxExchange = ContextWithExec(nil, execCtxProxy)

func execCtxProxy(req *Request) (tag CommandTag, err error) {
	switch req.Uri {
	case BuildUpdateUri(execTestUpdateRsc):
		return tag, errors.New("exec error")
	case BuildInsertUri(execTestInsertRsc):
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

func ExampleExecContext_Error() {
	req := NewUpdateRequest(execTestUpdateRsc, execTestUpdateSql, nil, nil)
	tag, err := Exec[runtime.DebugError](context.Background(), NullCount, req)
	fmt.Printf("test: Exec[DebugError](ctx,NullCount,update) : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//[[] github.com/gotemplates/postgresql/pgxsql/exec [error on PostgreSQL exec call : dbClient is nil]]
	//test: Exec[DebugError](ctx,NullCount,update) : [tag:{ 0 false false false false}] [error:InvalidArgument]

}

func ExampleExecContext_Insert() {
	//ctx := NewExecContext(nil, execCtxExchange)
	req := NewInsertRequest(execTestInsertRsc, execTestInsertSql, nil)
	tag, err := Exec[runtime.DebugError](execCtxExchange, NullCount, req)
	fmt.Printf("test: Exec[DebugError](ctx,NullCount,insert) : [tag:%v] [error:%v]\n", tag, err)

	//Output:
	//test: Exec[DebugError](ctx,NullCount,insert) : [tag:{INSERT 1 1234 true false false false}] [error:OK]

}

func ExampleExecContext() {
	k1 := "1"
	k2 := "2"
	v1 := "value 1"
	v2 := "value 2"

	ctx := execCtxExchange

	fmt.Printf("test: IsExecContext(ctx) -> %v\n", isExecContext(ctx))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx.Value(k1), ctx.Value(k2))

	ctx1 := ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: IsExecContext(ctx1) -> %v\n", isExecContext(ctx1))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx1.Value(k1), ctx1.Value(k2))

	ctx2 := ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: IsExecContext(ctx2) -> %v\n", isExecContext(ctx2))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx2.Value(k1), ctx2.Value(k2))

	//Output:
	//test: IsExecContext(ctx) -> true
	//test: Values() -> [key1:<nil>] [key2:<nil>]
	//test: IsExecContext(ctx1) -> true
	//test: Values() -> [key1:value 1] [key2:<nil>]
	//test: IsExecContext(ctx2) -> true
	//test: Values() -> [key1:value 1] [key2:value 2]

}
