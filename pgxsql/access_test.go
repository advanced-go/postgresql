package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/module"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/io"
	"github.com/advanced-go/stdlib/json"
	"net/http"
)

const (
	accessJson = "file://[cwd]/pgxsqltest/access.json"
)

func ExampleAccessInsert() {
	req := new(request)

	tag, err := accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(list))

	req.values = toValues(list)
	tag, err = accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(list))

	//Output:
	//test: accessInsert() -> [tag:{ 0 false false false false}] [err:request or request values is nil] [count:2]
	//test: accessInsert() -> [tag:{ 2 true false false false}] [err:<nil>] [count:4]

}

func ExampleInsert() {
	buf, status := io.ReadFile(accessJson)
	if !status.OK() {
		fmt.Printf("test: io.ReadFile() -> [status:%v]\n", status)
	} else {
		rows, status1 := json.New[[]Entry](buf, nil)
		if !status1.OK() {
			fmt.Printf("test: json.New() -> [status:%v]\n", status1)
		} else {
			values := toValues(rows)
			h := make(http.Header)
			h.Set(core.XFrom, module.Authority)
			tag, status2 := Insert(nil, h, "access-log", "", values)
			fmt.Printf("test: Insert() -> [tag:%v] [status:%v] [count:%v]\n", tag, status2, len(list))
		}
	}

	//Output:
	//test: Insert() -> [tag:{ 2 true false false false}] [status:OK] [count:6]

}

func toValues(entries []Entry) [][]any {
	var values [][]any
	for _, e := range entries {
		row := e.Values()
		values = append(values, row)
	}
	return values
}
