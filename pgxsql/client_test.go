package pgxsql

import (
	"fmt"
	"github.com/gotemplates/host/messaging"
)

func ExampleClientStartup() {

	db := messaging.DatabaseUrl{Url: ""}
	err := ClientStartup(db, nil)
	if err == nil {
		defer ClientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> database URL is empty

}
