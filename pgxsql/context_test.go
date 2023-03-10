package pgxsql

import (
	"context"
	"fmt"
)

func ExampleContextWithValue() {
	k1 := "1"
	v1 := "value 1"
	k2 := "2"
	v2 := "value 2"

	ctx := ContextWithValue(context.Background(), k1, v1)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	ctx = ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	//Output:
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:<nil>] [query:false] [exec:false]
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:value 2] [query:false] [exec:false]

}

func ExampleContextWithValue_Query() {
	k1 := "1"
	v1 := "value 1"
	k2 := "2"
	v2 := "value 2"

	ctx := ContextWithQuery(nil, nil)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	ctx = ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	ctx = ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	//Output:
	//test: ContextWithValue(ctx,k1,v1) -> [v1:<nil>] [v2:<nil>] [query:true] [exec:false]
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:<nil>] [query:true] [exec:false]
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:value 2] [query:true] [exec:false]

}

func ExampleContextWithValue_Exec() {
	k1 := "1"
	v1 := "value 1"
	k2 := "2"
	v2 := "value 2"

	ctx := ContextWithExec(nil, nil)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	ctx = ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	ctx = ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: ContextWithValue(ctx,k1,v1) -> [v1:%v] [v2:%v] [query:%v] [exec:%v]\n", ctx.Value(k1), ctx.Value(k2), isQueryContext(ctx), isExecContext(ctx))

	//Output:
	//test: ContextWithValue(ctx,k1,v1) -> [v1:<nil>] [v2:<nil>] [query:false] [exec:true]
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:<nil>] [query:false] [exec:true]
	//test: ContextWithValue(ctx,k1,v1) -> [v1:value 1] [v2:value 2] [query:false] [exec:true]

}
