package pgxsql

import "fmt"

func Example_QueryController() {
	qc := NewQueryController("query", Threshold{}, nil)

	fmt.Printf("test: NewQueryController() -> %v\n", qc)

	//Output:

}
