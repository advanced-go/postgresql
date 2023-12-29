package pgxdml

import (
	"fmt"
	"time"
)

func _ExampleFmtTimestamp() {
	t := time.Now().UTC()
	s := FmtTimestamp(t)
	fmt.Printf("test: FmtTimestamp() -> [%v]\n", s)

	t2, err := ParseTimestamp(s)
	fmt.Printf("test: ParseTimestamp() -> [%v] [%v]\n", FmtTimestamp(t2), err)

	//Output:

}

func stringer() string {
	s := "in stringer()"
	fmt.Printf("%v\n", s)
	return s
}

type stringerFunc func() string

func (f stringerFunc) String() string {
	return f()
}

func Example_StringerFunc() {
	var str fmt.Stringer

	str = stringerFunc(stringer)
	s := str.String()
	fmt.Printf("test: stringerFunc() -> %v\n", s)

	//Output:
	//in stringer()
	//test: stringerFunc() -> in stringer()

}
