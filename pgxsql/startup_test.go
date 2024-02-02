package pgxsql

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/uri"
	"net/http"
	"time"
)

// "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"

const (
	serviceUrl  = ""
	postgresUri = "github.com/idiomatic-go/postgresql/pgxsql"
)

func ExampleStartupPing() {
	r, _ := http.NewRequest("", "github/advanced-go/postgresql/pgxsql:ping", nil)
	nid, rsc, ok := uri.UprootUrn(r.URL.Path)
	status := messaging.Ping(nil, nid)
	fmt.Printf("test: Ping() -> [nid:%v] [nss:%v] [ok:%v] [status-code:%v]\n", nid, rsc, ok, status.Code)

	//Output:
	//test: Ping() -> [nid:github/advanced-go/postgresql/pgxsql] [nss:ping] [ok:true] [status-code:200]

}

func ExampleStartup() {
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
	//test: testStartup() -> [error:error running testStartup(): service url is empty]

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
