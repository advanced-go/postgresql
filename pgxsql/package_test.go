package pgxsql

const (
	configMapUri = "file://[cwd]/test/config-map.txt"
)

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
