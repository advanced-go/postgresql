package pgxdml

import (
	"fmt"
	"net/url"
	"strings"
)

func ExampleBuildWhere() {
	u, _ := url.Parse("http://www.google.com/search?loc=texas&zone=frisco")
	where := BuildWhere(u.Query())
	fmt.Printf("test: BuildWhere(u) -> %v\n", where)

	//Output:
	//test: BuildWhere(u) -> [{loc texas} {zone frisco}]

}

func ExampleWriteWhere() {
	sb := strings.Builder{}

	err := WriteWhere(&sb, false, nil)
	fmt.Printf("test: WriteWhere(false,nil) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	err = WriteWhere(&sb, false, []Attr{{Key: "", Val: nil}})
	fmt.Printf("test: WriteWhere(false,empty name) -> [error:%v] [stmt:%v]\n", err, NilEmpty(strings.Trim(sb.String(), " ")))

	sb.Reset()
	err = WriteWhere(&sb, true, []Attr{{Key: "status_code", Val: "503"}})
	fmt.Printf("test: WriteWhere(true,name,val) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	sb.Reset()
	err = WriteWhere(&sb, false, []Attr{{Key: "status_code", Val: "503"}, {Key: "minimum_code", Val: 99}, {Key: "created_ts", Val: Function("now()")}})
	fmt.Printf("test: WriteWhere(false,name value) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	//Output:
	//test: WriteWhere(false,nil) -> [error:invalid update where argument, attrs slice is empty] [stmt:<nil>]
	//test: WriteWhere(false,empty name) -> [error:<nil>] [stmt:WHERE]
	//test: WriteWhere(true,name,val) -> [error:<nil>] [stmt:WHERE status_code = '503';]
	//test: WriteWhere(false,name value) -> [error:<nil>] [stmt:WHERE status_code = '503' AND minimum_code = 99 AND created_ts = now()]

}

func ExampleWriteWhereAttributes() {
	sb := strings.Builder{}

	err := WriteWhereAttributes(&sb, nil)
	fmt.Printf("test: WriteWhereAttributes(nil) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	err = WriteWhereAttributes(&sb, []Attr{{Key: "", Val: nil}})
	fmt.Printf("test: WriteWhereAttributes(empty name) -> [error:%v] [stmt:%v]\n", err, NilEmpty(strings.Trim(sb.String(), " ")))

	sb.Reset()
	err = WriteWhereAttributes(&sb, []Attr{{Key: "status_code", Val: "503"}})
	fmt.Printf("test: WriteWhereAttributes(name,val) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	sb.Reset()
	err = WriteWhereAttributes(&sb, []Attr{{Key: "status_code", Val: "503"}, {Key: "minimum_code", Val: 99}, {Key: "created_ts", Val: Function("now()")}})
	fmt.Printf("test: WriteWhereAttributes(name value) -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	//Output:
	//test: WriteWhereAttributes(nil) -> [error:invalid update where argument, attrs slice is empty] [stmt:<nil>]
	//test: WriteWhereAttributes(empty name) -> [error:invalid attribute argument, attribute name is empty] [stmt:<nil>]
	//test: WriteWhereAttributes(name,val) -> [error:<nil>] [stmt:status_code = '503']
	//test: WriteWhereAttributes(name value) -> [error:<nil>] [stmt:status_code = '503' AND minimum_code = 99 AND created_ts = now()]

}
