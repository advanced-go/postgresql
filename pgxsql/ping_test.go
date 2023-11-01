package pgxsql

import (
	"fmt"
)

func ExamplePing() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isStarted())

		status := Ping(nil)
		fmt.Printf("test: Ping(nil) -> %v\n", status)
	}

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Ping(nil) -> OK

}
