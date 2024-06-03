package pgxsql

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/http"
)

const (
	configMapUri = "file://[cwd]/pgxsqltest/config-map.txt"
)

func ExampleQuery() {
	h := make(http.Header)
	h.Add(core.XAuthority, "github/advanced-go/timeseries")
	rows, status := Query(nil, h, "test-access.dev", "", nil)
	if !status.OK() {
		fmt.Printf("test: Query() -> [status:%v]\n", status)
	} else {
		entries, status1 := Scan[Entry](rows)
		fmt.Printf("test: Query() -> [status:%v] [entries:%v]\n", status1, len(entries))
	}

	//Output:
	//test: Query() -> [status:OK] [entries:2]

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
