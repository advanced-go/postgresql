package pgxsql

import (
	"fmt"
	"time"
)

func Example_NewStatusAgent() {
	a, err := NewStatusAgent(10, time.Second*2, nil, nil)

	fmt.Printf("test: NewStatusAgent() -> %v [err:%v]\n", a, err)

	//Output:
	//test: NewStatusAgent() -> &{2000000000 86400000000000 10 0xc00007e360 <nil> <nil>} [err:<nil>]

}
