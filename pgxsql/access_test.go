package pgxsql

import (
	"fmt"
	"github.com/advanced-go/postgresql/module"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/json"
	"github.com/advanced-go/stdlib/uri"
	"net/http"
	"net/url"
)

const (
	accessJson = "file://[cwd]/pgxsqltest/access.json"
)

func _ExampleAccessInsert() {
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

func _ExampleInsert() {
	rows, status := json.New[[]Entry](accessJson, nil)
	if !status.OK() {
		fmt.Printf("test: io.ReadFile() -> [status:%v]\n", status)
	} else {
		values := toValues(rows)
		h := make(http.Header)
		h.Set(core.XFrom, module.Authority)
		tag, status2 := Insert(nil, h, "access-log", "", values)
		fmt.Printf("test: Insert() -> [tag:%v] [status:%v] [count:%v]\n", tag, status2, len(list))
	}

	//Output:
	//test: Insert() -> [tag:{ 2 true false false false}] [status:OK] [count:6]

}

func _ExampleQueryT() {
	h := make(http.Header)
	h.Set(core.XFrom, module.Authority)
	values := make(url.Values)
	values.Add(core.RegionKey, "*")
	entries, status := QueryT[Entry](nil, h, "access-log", "", values)

	fmt.Printf("test: QueryT() -> [status:%v] [entries:%v]\n", status, len(entries))

	//Output:
	//test: QueryT() -> [status:OK] [entries:6]

}

func toValues(entries []Entry) [][]any {
	var values [][]any
	for _, e := range entries {
		row := e.Values()
		values = append(values, row)
	}
	return values
}

func _ExampleAccessFilter() {
	q := ""
	result := accessFilter(nil)
	fmt.Printf("test: accessFilter(\"%v\") -> [cnt:%v] [filter:%v]\n", q, len(list), len(result))

	q = "region=*&order=desc"
	result = accessFilter(uri.BuildValues(q))
	fmt.Printf("test: accessFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(list), len(result), nil)

	q = "region=*&order=desc&top=2"
	result = accessFilter(uri.BuildValues(q))
	fmt.Printf("test: accessFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(list), len(result), nil)

	q = "region=*&order=desc&top=45"
	result = accessFilter(uri.BuildValues(q))
	fmt.Printf("test: accessFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(list), len(result), nil)

	//Output:
	//test: accessFilter("") -> [cnt:4] [filter:4]
	//test: accessFilter("region=*&order=desc") -> [cnt:4] [filter:4] [result:<nil>]
	//test: accessFilter("region=*&order=desc&top=2") -> [cnt:4] [filter:2] [result:<nil>]
	//test: accessFilter("region=*&order=desc&top=45") -> [cnt:4] [filter:4] [result:<nil>]

}

func _ExampleOrder() {
	q := ""
	result := order(nil, list)
	fmt.Printf("test: order(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), result)

	q = "order=desc"
	result = order(uri.BuildValues(q), list)
	fmt.Printf("test: order(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), result)

	//Output:
	//fail
}

func ExampleTop() {
	q := ""
	result := top(nil, list)
	fmt.Printf("test: top(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), len(result))

	q = "top=2"
	result = top(uri.BuildValues(q), list)
	fmt.Printf("test: top(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), len(result))

	//Output:
	//test: top("") -> [cnt:4] [result:4]
	//test: top("top=2") -> [cnt:4] [result:2]

}

func ExampleDistinct() {
	q := ""
	result := distinct(nil, list)
	fmt.Printf("test: distinct(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), len(result))

	q = "distinct=host"
	result = distinct(uri.BuildValues(q), list)
	fmt.Printf("test: distinct(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(list), len(result))

	//Output:
	//test: distinct("") -> [cnt:4] [result:4]
	//test: distinct("distinct=host") -> [cnt:4] [result:3]

}
