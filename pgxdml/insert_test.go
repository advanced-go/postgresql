package pgxdml

import (
	"fmt"
	"strings"
)

const (
	insertEntryStmt = "INSERT INTO test_entry (id,customer_id,ping_traffic,counter_value,changed_ts) VALUES"
	//nextValFn       = "nextval('test_entry_Id')"
)

func NilEmpty(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

func ExampleNewInsertValues() {
	values := NewInsertValues([]any{100, "customer 1", false, NextValFn("test_entry_Id"), TimestampFn})

	fmt.Printf("test: NewInsertValues() -> %v\n", values)

	//Output:
	//test: NewInsertValues() -> [[100 customer 1 false nextval('test_entry_Id') now()]]
}

func ExampleWriteInsert() {
	var values [][]any
	values = append(values, []any{100, "customer 1", false, NextValFn("test_entry_Id"), TimestampFn})
	values = append(values, []any{200, "customer 2", true, NextValFn("test_entry_Id"), TimestampFn})

	stmt, err := WriteInsert(insertEntryStmt, values)
	fmt.Printf("test: WriteInsert() -> [error:%v] [stmt:%v\n", err, stmt)

	//Output:
	//test: WriteInsert() -> [error:<nil>] [stmt:INSERT INTO test_entry (id,customer_id,ping_traffic,counter_value,changed_ts) VALUES
	//(100,'customer 1',false,nextval('test_entry_Id'),now()),
	//(200,'customer 2',true,nextval('test_entry_Id'),now());

}

func ExampleWriteInsertValues() {
	sb := strings.Builder{}

	err := WriteInsertValues(&sb, nil)
	fmt.Printf("test: WriteInsertValues() -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	sb1 := strings.Builder{}
	err = WriteInsertValues(&sb1, []any{100})
	fmt.Printf("test: WriteInsertValues() -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb1.String()))

	err = WriteInsertValues(&sb, []any{100, "test string", false, NextValFn("test_entry_Id"), TimestampFn})
	fmt.Printf("test: WriteInsertValues() -> [error:%v] [stmt:%v]\n", err, NilEmpty(sb.String()))

	//Output:
	//test: WriteInsertValues() -> [error:invalid insert argument, values slice is empty] [stmt:<nil>]
	//test: WriteInsertValues() -> [error:<nil>] [stmt:(100)]
	//test: WriteInsertValues() -> [error:<nil>] [stmt:(100,'test string',false,nextval('test_entry_Id'),now())]

}
