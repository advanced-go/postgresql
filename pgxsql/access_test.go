package pgxsql

import "fmt"

func ExampleAccessInsert() {
	req := new(request)

	tag, err := accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(list))

	for _, e := range list {
		row := e.Values()
		req.values = append(req.values, row)
	}
	tag, err = accessInsert(nil, "", req)
	fmt.Printf("test: accessInsert() -> [tag:%v] [err:%v] [count:%v]\n", tag, err, len(list))

	//Output:
	//test: accessInsert() -> [tag:{ 0 false false false false}] [err:request or request values is nil] [count:2]
	//test: accessInsert() -> [tag:{ 2 true false false false}] [err:<nil>] [count:4]

}
