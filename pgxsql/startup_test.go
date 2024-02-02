package pgxsql

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"time"
)

// "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"

const (
	serviceUrl  = ""
	postgresUri = "github.com/idiomatic-go/postgresql/pgxsql"
)

func Example_Startup() {
	fmt.Printf("test: isReady() -> %v\n", isReady())
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isReady())

		//status := host.Ping[runtime.Output](nil, postgresUri)
		//fmt.Printf("test: messaging.Ping() -> %v\n", status)

	}

	//Output:
	//test: isReady() -> false
	//test: clientStartup() -> [started:true]
	//{traffic:egress, route:*, request-id:, status-code:0, method:GET, url:urn:postgres:ping, startup:postgres, path:ping, timeout:-1, rate-limit:-1, rate-burst:-1, retry:, retry-rate-limit:-1, retry-rate-burst:-1, status-flags:}
	//test: messaging.Ping() -> OK

}

func testStartup() error {
	if serviceUrl == "" {
		return errors.New("error running testStartup(): service url is empty")
	}
	if isReady() {
		return nil
	}

	m := make(map[string]string)
	m[uriConfigKey] = serviceUrl
	messaging.HostExchange.SendCtrl(messaging.Message{
		To:    PkgPath,
		From:  "",
		Event: messaging.StartupEvent,
		//Status:  nil,
		Config:  m,
		ReplyTo: nil,
	})
	time.Sleep(time.Second * 3)
	return nil
}
