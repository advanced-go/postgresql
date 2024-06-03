package pgxsql

import (
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/host"
	"github.com/advanced-go/stdlib/messaging"
	"net/http"
	"time"
)

// "postgres://{user}:{pswd}@{sub-domain}.{database}.cloud.timescale.com:{port}/{database}?sslmode=require"

const (
	serviceUrl = ""
)

func ExampleStartupPing() {
	r, _ := http.NewRequest("", PkgPath+":ping", nil)
	status := messaging.Ping(nil, r.URL)
	fmt.Printf("test: Ping() -> [status-code:%v]\n", status.Code)

	//Output:
	//test: Ping() -> [status-code:200]

}

func ExampleStartup() {
	fmt.Printf("test: isReady() -> %v\n", isReady())
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		fmt.Printf("test: clientStartup() -> [started:%v]\n", isReady())

		//status := host.Ping[core.Output](nil, postgresUri)
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
	msg := messaging.NewMessage(messaging.ChannelControl, PkgPath, "", messaging.StartupEvent)
	msg.SetContent(messaging.ContentTypeConfig, m)
	host.Exchange.Send(msg)
	time.Sleep(time.Second * 3)
	return nil
}
