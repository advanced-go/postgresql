package pgxsql

import (
	"fmt"
)

func ExampleStat() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isStarted())

		stat, status := Stat(nil)
		fmt.Printf("test: Stat(nil) -> [status:%v] [stat:%v]\n", status, stat != nil)
	}

	//Output:
	//test: clientStartup() -> [started:true]
	//test: Stat(nil) -> [status:OK] [stat:true]

}
