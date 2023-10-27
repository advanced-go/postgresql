package pgxdml

import (
	"fmt"
	"time"
)

/*
func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}


*/

func ExampleFmtValues() {
	var ptr *int

	v, err := FmtValue(nil)
	fmt.Printf("test: FmtValue(nil) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue(ptr)
	fmt.Printf("test: FmtValue(ptr) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	var n = 123
	v, err = FmtValue(&n)
	fmt.Printf("test: FmtValue(&n) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue(true)
	fmt.Printf("test: FmtValue(true) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue(1001)
	fmt.Printf("test: FmtValue(1001) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue("")
	fmt.Printf("test: FmtValue(\"\") -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	//v, err = FmtValue(NowFn)
	//fmt.Printf("test: FmtValue(now()) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue(time.Date(2023, 2, 5, 13, 25, 56, 0, time.UTC))
	fmt.Printf("test: FmtValue(time.Now()) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue("test string")
	fmt.Printf("test: FmtValue(test string) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	//Output:
	//test: FmtValue(nil) -> [error:<nil>] [value:NULL]
	//test: FmtValue(ptr) -> [error:<nil>] [value:NULL]
	//test: FmtValue(&n) -> [error:invalid argument : pointer types are not supported : *int] [value:<nil>]
	//test: FmtValue(true) -> [error:<nil>] [value:true]
	//test: FmtValue(1001) -> [error:<nil>] [value:1001]
	//test: FmtValue("") -> [error:<nil>] [value:'']
	//test: FmtValue(time.Now()) -> [error:<nil>] [value:'2023-02-05 13:25:56.000000']
	//test: FmtValue(test string) -> [error:<nil>] [value:'test string']

}

func ExampleFmtSqlValues() {
	v, err := FmtValue(Function("now()"))
	fmt.Printf("test: FmtValue(now()) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	v, err = FmtValue("drop table")
	fmt.Printf("test: FmtValue(drop table) -> [error:%v] [value:%v]\n", err, NilEmpty(v))

	//Output:
	//test: FmtValue(now()) -> [error:<nil>] [value:now()]
	//test: FmtValue(drop table) -> [error:SQL injection embedded in string [drop table] : drop table] [value:<nil>]

}

func ExampleFmtAttr() {
	s, err := FmtAttr(Attr{})
	fmt.Printf("Name  [\"\"]  : %v\n", NilEmpty(s))
	fmt.Printf("Error       : %v\n", err)

	s, err = FmtAttr(Attr{Key: "attr_name_1"})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Key: "attr_name_2", Val: 1234})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Key: "attr_name_3", Val: false})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	//s, err = FmtAttr(util.Attr{Name: "attr_name_4", Val: time.Now()})
	//fmt.Println("default format:", time.Now())
	//fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	//fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Key: "attr_name_5", Val: "value string"})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	s, err = FmtAttr(Attr{Key: "attr_name_6", Val: Function("now()")})
	fmt.Printf("Name  [attr_name]  : %v\n", NilEmpty(s))
	fmt.Printf("Error              : %v\n", err)

	//Output:
	//Name  [""]  : <nil>
	//Error       : invalid attribute argument, attribute name is empty
	//Name  [attr_name]  : attr_name_1 = NULL
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_2 = 1234
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_3 = false
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_5 = 'value string'
	//Error              : <nil>
	//Name  [attr_name]  : attr_name_6 = now()
	//Error              : <nil>

}
