package pgxsql

import "fmt"

func ExampleAccessRows() {
	rows, _ := accessQuery(nil, "")
	entries, status := Scan[Entry](rows)
	fmt.Printf("test: accessRows() -> [status:%v] [entries:%v]\n", status, len(entries)) //entries)

	//Output:
	//test: accessRows() -> [status:OK] [entries:2]

}
