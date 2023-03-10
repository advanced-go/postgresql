package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/gotemplates/core/runtime"
)

const (
	queryTestErrorSql = "select * from test"
	queryTestRowsSql  = "select * from table"
	queryTestErrorRsc = "error"
	queryTestRowsRsc  = "rows"
)

func isQueryContext(ctx context.Context) bool {
	if _, ok := any(ctx).(queryWithValue); ok {
		return true
	}
	return false
}

var queryCtxExchange = ContextWithQuery(nil, queryCtxProxy)

func queryCtxProxy(req *Request) (Rows, error) {
	switch req.Uri {
	case BuildQueryUri(queryTestErrorRsc):
		return nil, errors.New("pgxsql query error")
	case BuildQueryUri(queryTestRowsRsc):
		var i Rows = &rowsT{}
		return i, nil
	}
	return nil, nil
}

func ExampleQueryContext_Error() {
	ctx := queryCtxExchange
	req := NewQueryRequest(queryTestErrorRsc, queryTestErrorSql, nil)
	rows, status := Query[runtime.DebugError](ctx, req)
	fmt.Printf("test: Query[DebugError](ctx,req) : [rows:%v] [status:%v]\n", rows != nil, status)

	//Output:
	//[[] github.com/gotemplates/postgresql/pgxsql/exec [pgxsql query error]]
	//test: Query[DebugError](ctx,req) : [rows:false] [status:Internal]

}

func ExampleQueryContext_Rows() {
	ctx := queryCtxExchange
	req := NewQueryRequest(queryTestRowsRsc, queryTestRowsSql, nil)
	rows, status := Query[runtime.DebugError](ctx, req)
	fmt.Printf("test: Query[DebugError](ctx,req) : [rows:%v] [status:%v]\n", rows != nil, status)

	//Output:
	//test: Query[DebugError](ctx,req) : [rows:true] [status:OK]

}

func ExampleQueryContext() {
	k1 := "1"
	k2 := "2"
	v1 := "value 1"
	v2 := "value 2"

	ctx := queryCtxExchange

	fmt.Printf("test: IsQueryContext(ctx) -> %v\n", isQueryContext(ctx))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx.Value(k1), ctx.Value(k2))

	ctx1 := ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: IsQueryContext(ctx1) -> %v\n", isQueryContext(ctx1))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx1.Value(k1), ctx1.Value(k2))

	ctx2 := ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: IsQueryContext(ctx2) -> %v\n", isQueryContext(ctx2))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v]\n", ctx2.Value(k1), ctx2.Value(k2))

	//Output:
	//test: IsQueryContext(ctx) -> true
	//test: Values() -> [key1:<nil>] [key2:<nil>]
	//test: IsQueryContext(ctx1) -> true
	//test: Values() -> [key1:value 1] [key2:<nil>]
	//test: IsQueryContext(ctx2) -> true
	//test: Values() -> [key1:value 1] [key2:value 2]

}
