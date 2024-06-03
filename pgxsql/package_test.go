package pgxsql

import (
	"fmt"
	"reflect"
)

const (
	configMapUri = "file://[cwd]/pgxsqltest/config-map.txt"
)

func _Example_PackageUri() {
	pkgPath := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath)

	//Output:
	//test: PkgPath = "github.com/advanced-go/postgresql/pgxsql"

}

/*
func ExampleNewStringsMap() {
	uri := configMapUri


	fmt.Printf("test: NewStringsMap(\"%v\") -> [err:%v]\n", uri, m.Error())

	key := userConfigKey
	val, status := m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [user:%v] [status:%v]\n", key, val, status)

	key = pswdConfigKey
	val, status = m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [pswd:%v] [status:%v]\n", key, val, status)

	key = uriConfigKey
	val, status = m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [urir:%v] [status:%v]\n", key, val, status)

	//Output:
	//test: NewStringsMap("file://[cwd]/pgxsqltest/config-map.txt") -> [err:<nil>]
	//test: Get("user") -> [user:bobs-your-uncle] [status:OK]
	//test: Get("pswd") -> [pswd:let-me-in] [status:OK]
	//test: Get("uri") -> [urir:postgres://{user}:{pswd}@{sub-domain}.{db-name}.cloud.timescale.com:31770/tsdb?sslmode=require] [status:OK]

}


*/
