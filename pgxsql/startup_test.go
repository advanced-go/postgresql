package pgxsql

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"github.com/advanced-go/messaging/exchange"
	"time"
)

// "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"

const (
	serviceUrl  = ""
	postgresUri = "github.com/idiomatic-go/postgresql/pgxsql"
)

func Example_Startup() {
	fmt.Printf("test: isStarted() -> %v\n", isStarted())
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer ClientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isStarted())

		status := exchange.Ping[runtime.TestError](nil, postgresUri)
		fmt.Printf("test: messaging.Ping() -> %v\n", status)

	}

	//Output:
	//test: isStarted() -> false
	//test: clientStartup() -> [started:true]
	//{traffic:egress, route:*, request-id:, status-code:0, method:GET, url:urn:postgres:ping, startup:postgres, path:ping, timeout:-1, rate-limit:-1, rate-burst:-1, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:}
	//test: messaging.Ping() -> OK

}

func testStartup() error {
	if serviceUrl == "" {
		return errors.New("error running testStartup(): service url is empty")
	}
	if isStarted() {
		return nil
	}

	c <- core.Message{
		To:      "",
		From:    "",
		Event:   core.StartupEvent,
		Status:  nil,
		Content: []any{core.Resource{Uri: serviceUrl}},
		ReplyTo: nil,
	}
	time.Sleep(time.Second * 3)

	return nil
}
