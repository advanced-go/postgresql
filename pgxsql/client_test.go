package pgxsql

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime/startup"
)

func ExampleClientStartup() {

	rsc := startup.Resource{Uri: ""}
	err := ClientStartup(rsc, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
