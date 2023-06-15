package pgxdml

import (
	"fmt"
	"github.com/go-sre/core/runtime"
)

const (
	deleteTestEntryStmt = "DELETE test_entry"
)

func ExampleWriteDelete() {
	where := []runtime.Attr{{Key: "customer_id", Val: "customer1"}, {Key: "created_ts", Val: "2022/11/30 15:48:54.049496"}} //time.Now()}}

	sql, err := WriteDelete(deleteTestEntryStmt, where)
	fmt.Printf("test: WriteDelete() -> [error:%v] [stmt:%v]\n", err, NilEmpty(sql))

	//Output:
	//test: WriteDelete() -> [error:<nil>] [stmt:DELETE test_entry
	//WHERE customer_id = 'customer1' AND created_ts = '2022/11/30 15:48:54.049496';]

}
