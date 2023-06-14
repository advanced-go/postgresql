package pgxdml

import "fmt"

const (
	// TimestampFn - timestamp SQL function
	TimestampFn = Function("now()")

	nextValFnFmt = "nextval('%s')"
	valueFmt     = "%v"
	stringFmt    = "'%v'"
	attrFmt      = "%v = %v"
)

// NextValFn - build a postgresql SQL 'nextval()' function
func NextValFn(sequence string) Function {
	return Function(fmt.Sprintf(nextValFnFmt, sequence))
}

// Function - type used to determine formatting of a functions
type Function string

var tokens = []string{"drop table", "delete from", "--", ";", "/*", "*/", "select * from"}
