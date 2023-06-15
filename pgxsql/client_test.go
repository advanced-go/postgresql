package pgxsql

import (
	"fmt"
	"github.com/go-ai-agent/core/resource"
)

func ExampleClientStartup() {

	db := resource.DatabaseUrl{Url: ""}
	err := ClientStartup(db, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
