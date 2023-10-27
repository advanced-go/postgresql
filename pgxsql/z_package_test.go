package pgxsql

import "fmt"

func Example_PackageUri() {
	fmt.Printf("test: pkgUri -> %v\n", pkgUri)
	fmt.Printf("test: pkgUri -> %v\n", pkgPath)

	//Output:
	//test: pkgUri -> github.com/go-ai-agent/postgresql/pgxsql
	//test: pkgUri -> /go-ai-agent/postgresql/pgxsql

}
