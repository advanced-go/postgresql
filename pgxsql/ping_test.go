package pgxsql

import (
	"fmt"
	"github.com/go-ai-agent/core/controller"
	"github.com/go-ai-agent/core/runtime"
)

func ExamplePing() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", IsStarted())

		status := Ping[runtime.DebugError, controller.NilHandler](nil)
		fmt.Printf("test: Ping(nil) -> %v\n", status)
	}

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Ping(nil) -> OK

}
