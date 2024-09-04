package access

import (
	"fmt"
	"github.com/advanced-go/stdlib/uri"
)

func _ExampleOrder() {
	q := ""
	result := order(nil, storage)
	fmt.Printf("test: order(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), result)

	q = "order=desc"
	result = order(uri.BuildValues(q), storage)
	fmt.Printf("test: order(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), result)

	//Output:
	//fail
}

func ExampleTop() {
	q := ""
	result := top(nil, storage)
	fmt.Printf("test: top(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), len(result))

	q = "top=2"
	result = top(uri.BuildValues(q), storage)
	fmt.Printf("test: top(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), len(result))

	//Output:
	//test: top("") -> [cnt:4] [result:4]
	//test: top("top=2") -> [cnt:4] [result:2]

}

func ExampleDistinct() {
	q := ""
	result := distinct(nil, storage)
	fmt.Printf("test: distinct(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), len(result))

	q = "distinct=host"
	result = distinct(uri.BuildValues(q), storage)
	fmt.Printf("test: distinct(\"%v\") -> [cnt:%v] [result:%v]\n", q, len(storage), len(result))

	//Output:
	//test: distinct("") -> [cnt:4] [result:4]
	//test: distinct("distinct=host") -> [cnt:4] [result:3]

}
