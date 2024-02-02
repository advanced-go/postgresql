package pgxsql

import (
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
	"sync/atomic"
	"time"
)

var (
	ready int64
	agent messaging.Agent
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
	var status runtime.Status
	agent, status = messaging.NewDefaultAgent(PkgPath, messageHandler, false)
	if !status.OK() {
		fmt.Printf("init(\"%v\") failure: [%v]\n", PkgPath, status)
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
		messaging.SendReply(msg, ping(nil).SetDuration(time.Since(start)))
	}
}
