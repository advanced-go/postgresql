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
	accessJson   = "file://[cwd]/pgxsqltest/access.json"
	accessStatus = "file://[cwd]/pgxsqltest/status-504.json"
)

func _ExampleAccessInsert() {
	req := new(request)

	tag, err := accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(storage))

	req.values = toValues(storage)
	tag, err = accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(storage))

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
		fmt.Printf("test: Insert() -> [tag:%v] [status:%v] [count:%v]\n", tag, status2, len(storage))
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

func ExampleQueryT_URL() {
	h := make(http.Header)
	h.Set(core.XFrom, module.Authority)
	values := make(url.Values)
	values.Add(core.RegionKey, "*")
	ctx := core.NewExchangeOverrideContext(nil, core.NewExchangeOverride("", accessJson, ""))
	entries, status := QueryT[Entry](ctx, h, "access-log", "", values)

	fmt.Printf("test: QueryT() -> [status:%v] [entries:%v]\n", status, len(entries))

	//Output:
	//test: QueryT() -> [status:OK] [entries:2]

}

func ExampleQueryT_URL_Status() {
	h := make(http.Header)
	h.Set(core.XFrom, module.Authority)
	values := make(url.Values)
	values.Add(core.RegionKey, "*")

	ctx := core.NewExchangeOverrideContext(nil, core.NewExchangeOverride("", "", accessStatus))
	entries, status := QueryT[Entry](ctx, h, "access-log", "", values)
	fmt.Printf("test: QueryT() -> [status:%v] [entries:%v]\n", status, len(entries))

	ctx = core.NewExchangeOverrideContext(nil, core.NewExchangeOverride("", "", json.StatusNotFoundUri))
	entries, status = QueryT[Entry](ctx, h, "access-log", "", values)
	fmt.Printf("test: QueryT() -> [status:%v] [entries:%v]\n", status, len(entries))

	//Output:
	//test: QueryT() -> [status:Timeout [error 1]] [entries:0]
	//test: QueryT() -> [status:Not Found] [entries:0]

}

func toValues(entries []Entry) [][]any {
	var values [][]any
	for _, e := range entries {
		row := e.Values()
		values = append(values, row)
	}
	return values
}

func ExampleOriginFilter() {
	q := ""
	result := originFilter(nil)
	fmt.Printf("test: originFilter(\"%v\") -> [cnt:%v] [filter:%v]\n", q, len(storage), len(result))

	q = "region=*&order=desc"
	result = originFilter(uri.BuildValues(q))
	fmt.Printf("test: originFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(storage), len(result), nil)

	q = "region=*&order=desc&top=2"
	result = originFilter(uri.BuildValues(q))
	fmt.Printf("test: originFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(storage), len(result), nil)

	q = "region=*&order=desc&top=45"
	result = originFilter(uri.BuildValues(q))
	fmt.Printf("test: originFilter(\"%v\") -> [cnt:%v] [filter:%v] [result:%v]\n", q, len(storage), len(result), nil)

	//Output:
	//test: originFilter("") -> [cnt:4] [filter:4]
	//test: originFilter("region=*&order=desc") -> [cnt:4] [filter:4] [result:<nil>]
	//test: originFilter("region=*&order=desc&top=2") -> [cnt:4] [filter:2] [result:<nil>]
	//test: originFilter("region=*&order=desc&top=45") -> [cnt:4] [filter:4] [result:<nil>]

}
