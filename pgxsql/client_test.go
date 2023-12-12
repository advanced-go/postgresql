package pgxsql

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
)

func ExampleClientStartup() {

	rsc := core.Resource{Uri: ""}
	err := clientStartup2(rsc, nil)
	if err == nil {
		defer clientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
