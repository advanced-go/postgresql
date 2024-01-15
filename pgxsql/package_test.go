package pgxsql

import (
	"fmt"
	"reflect"
)

func _Example_PackageUri() {
	pkgPath := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath)

	//Output:
	//test: PkgPath = "github.com/advanced-go/postgresql/pgxsql"

}
