package pgxsql

import "fmt"

func Example_QueryController() {
	qc := NewQueryController("query", Threshold{})

	fmt.Printf("test: NewQueryController() -> %v\n", qc)

	//Output:

}
