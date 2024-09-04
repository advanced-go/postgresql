package access

import (
	"fmt"
	"github.com/advanced-go/postgresql/pgxsql"
	"github.com/advanced-go/stdlib/uri"
)

func ExampleAccessQuery_All() {
	rows, _ := accessQuery(nil, "", new(pgxsql.request))
	entries, status := pgxsql.Scan[Entry](rows)
	fmt.Printf("test: accessQuery() -> [status:%v] [entries:%v]\n", status, len(entries)) //entries)

	//Output:
	//test: accessQuery() -> [status:OK] [entries:4]

}

func ExampleAccessQuery_Distinct() {
	req := new(pgxsql.request)
	req.values2 = uri.BuildValues("region=*&distinct=host")
	rows, _ := accessQuery(nil, "", req)
	entries, status := pgxsql.Scan[Entry](rows)
	fmt.Printf("test: accessQuery() -> [status:%v] [entries:%v]\n", status, len(entries)) //entries)

	//Output:
	//test: accessQuery() -> [status:OK] [entries:3]

}
