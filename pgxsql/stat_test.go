package pgxsql

import (
	"fmt"
)

func ExampleStat() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isReady())

		stat1, status := stat()
		fmt.Printf("test: Stat() -> [status:%v] [stat:%v]\n", status, stat1 != nil)
	}

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Stat(nil) -> [status:OK] [stat:true]

}
