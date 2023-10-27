package pgxdml

import (
	"fmt"
)

func ExampleExpandSelect() {
	t := "select * from access_log {where} order by start_time desc limit 5"
	where := []Attr{{Key: "status_code", Val: "503"}}

	sql, err := ExpandSelect("", nil)
	fmt.Printf("test: ExpandSelect(nil,nil) -> [error:%v] [empty:%v]\n", err, sql == "")

	sql, err = ExpandSelect(t, nil)
	fmt.Printf("test: ExpandSelect(t,nil) -> [error:%v] %v\n", err, sql)

	sql, err = ExpandSelect(t, where)
	fmt.Printf("test: ExpandSelect(t,where) -> [error:%v] %v\n", err, sql)

	//Output:
	//test: ExpandSelect(nil,nil) -> [error:template is empty] [empty:true]
	//test: ExpandSelect(t,nil) -> [error:<nil>] select * from access_log order by start_time desc limit 5
	//test: ExpandSelect(t,where) -> [error:<nil>] select * from access_log
	//WHERE status_code = '503' order by start_time desc limit 5

}
