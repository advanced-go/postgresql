package pgxsql

import (
	"fmt"
	"github.com/advanced-go/core/messaging"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	ready int64
	agent *messaging.Agent
)

func isReady() bool {
	return atomic.LoadInt64(&ready) != 0
}

func setReady() {
	atomic.StoreInt64(&ready, 1)
}

func resetReady() {
	atomic.StoreInt64(&ready, 0)
}

func init() {
	var err error
	agent, err = messaging.NewDefaultAgent(PkgPath, messageHandler, false)
	if err != nil {
		fmt.Printf("init(\"%v\") failure: [%v]\n", PkgPath, err)
	}
	agent.Run()
}

var messageHandler messaging.MessageHandler = func(msg messaging.Message) {
	switch msg.Event {
	case messaging.StartupEvent:
		clientStartup(msg)
	case messaging.ShutdownEvent:
		clientShutdown()
	case messaging.PingEvent:
		start := time.Now()
		messaging.SendReply(msg, messaging.NewStatusDuration(http.StatusOK, time.Since(start)))
	}
}
