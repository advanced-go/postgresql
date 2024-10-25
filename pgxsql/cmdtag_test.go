package pgxsql

import (
	"fmt"
)

var (
	//var cmdJson = "{\"sql\":\"select * from table\",\"rows-affected\":123,\"insert\":false,\"update\":true,\"delete\":false,\"select\":false}"
	updateCmdTag2 = "file://[cwd]/test/update-cmd-tag.json"
)

func ExampleNewCommandTag() {
	/*
		cmd := CommandTag{
			Sql:          "select * from table",
			RowsAffected: 123,
			Insert:       false,
			Update:       true,
			Delete:       false,
			Select:       false,
		}
		buf, status := json.Marshal(cmd)
		fmt.Printf("test: NewCommandTag() -> %v [status:%v]\n", string(buf), status)

	*/

	cmd2 := NewCommandTag(updateCmdTag2)
	fmt.Printf("test: NewCommandTag() -> %v\n", cmd2)

	//Output:
	//test: NewCommandTag() -> {update test 123 false true false false}

}
