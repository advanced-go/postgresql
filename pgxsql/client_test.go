package pgxsql

import (
	"fmt"
)

func ExampleClientStartup() {

	rsc := StartupResource{Uri: ""}
	err := clientStartup2(rsc, nil)
	if err == nil {
		defer clientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
