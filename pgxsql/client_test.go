package pgxsql

import (
	"fmt"
	"github.com/advanced-go/messaging/content"
)

func ExampleClientStartup() {

	rsc := content.Resource{Uri: ""}
	err := ClientStartup(rsc, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
