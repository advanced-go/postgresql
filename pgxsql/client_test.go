package pgxsql

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

func ExampleClientStartup() {
	//rsc := startupResource{Uri: ""}
	err := clientStartup2(nil)
	if err != nil {
		defer clientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	err = clientStartup2(runtime.NewEmptyStringsMap())
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> error: strings map configuration is nil
	//test: ClientStartup() -> database URL is empty

}
