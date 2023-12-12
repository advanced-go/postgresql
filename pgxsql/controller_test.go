package pgxsql

import "fmt"

func Example_QueryController() {
	qc := newQueryController("query", thresholdValues{}, nil)

	fmt.Printf("test: NewQueryController() -> %v\n", qc)

	//Output:
	//test: NewQueryController() -> &{query {0 0 0} <nil> <nil>}

}
