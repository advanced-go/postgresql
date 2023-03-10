package pgxdml

import "fmt"

func ExampleSanitizeString() {
	err := SanitizeString("")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("adfsdfe4fc&*4")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 1: /*")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 2: DROP   Table  ")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 3: DEL ETE FROM")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 4: - -")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 5: ;--")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 6: sa*/nitize 4 ;--")
	fmt.Printf("Error  : %v\n", err)

	err = SanitizeString("test 7: of select  * froM customers")
	fmt.Printf("Error  : %v\n", err)

	//Output:
	//Error  : <nil>
	//Error  : <nil>
	//Error  : SQL injection embedded in string [test 1: /*] : /*
	//Error  : SQL injection embedded in string [test 2: drop table] : drop table
	//Error  : <nil>
	//Error  : <nil>
	//Error  : SQL injection embedded in string [test 5: ;--] : --
	//Error  : SQL injection embedded in string [test 6: sa*/nitize 4 ;--] : --
	//Error  : SQL injection embedded in string [test 7: of select * from customers] : select * from
}
