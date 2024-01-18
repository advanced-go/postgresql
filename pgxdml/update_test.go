package pgxdml

import (
	"fmt"
	"strings"
)

const (
	updateTestEntryStmt = "UPDATE test_entry"
)

func ExampleWriteUpdate() {
	where := []Attr{{Key: "customer_id", Val: "customer1"}, {Key: "created_ts", Val: "2022/11/30 15:48:54.049496"}} //time.Now()}}
	attrs := []Attr{{Key: "status_code", Val: "503"}, {Key: "minimum_code", Val: 99}, {Key: "created_ts", Val: Function("now()")}}

	sql, err := WriteUpdate(updateTestEntryStmt, attrs, where)
	fmt.Printf("test: WriteUpdate(stmt,attrs,where) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sql))

	//fmt.Printf("Stmt       : %v\n", NilEmpty(sql))
	//fmt.Printf("Error      : %v\n", err)

	//Output:
	//test: WriteUpdate(stmt,attrs,where) -> [error:<nil>] [stmt:UPDATE test_entry
	//SET status_code = '503',
	//minimum_code = 99,
	//created_ts = now()
	//WHERE customer_id = 'customer1' AND created_ts = '2022/11/30 15:48:54.049496';]

}

func ExampleWriteUpdateSet() {
	sb := strings.Builder{}

	err := WriteUpdateSet(&sb, nil)
	fmt.Printf("test: WriteUpdateSet(nil) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	sb.Reset()
	err = WriteUpdateSet(&sb, []Attr{{Key: "status_code", Val: "503"}})
	fmt.Printf("test: WriteUpdateSet(name value) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	sb.Reset()
	err = WriteUpdateSet(&sb, []Attr{{Key: "status_code", Val: "503"}, {Key: "minimum_code", Val: 99}, {Key: "created_ts", Val: Function("now()")}})
	fmt.Printf("test: WriteUpdateSet(name value) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	//Output:
	//test: WriteUpdateSet(nil) -> [error:invalid update set argument, attrs slice is empty] [stmt:<nil>]
	//test: WriteUpdateSet(name value) -> [error:<nil>] [stmt:SET status_code = '503'
	//]
	//test: WriteUpdateSet(name value) -> [error:<nil>] [stmt:SET status_code = '503',
	//minimum_code = 99,
	//created_ts = now()
	//]

}
