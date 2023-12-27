package pgxsql

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/messaging/core"
	"github.com/advanced-go/messaging/exchange"
	"sync/atomic"
	"time"
)

var (
	pingThreshold   = 500
	queryThreshold  = 2000
	insertThreshold = 2000
	updateThreshold = 2000
	deleteThreshold = 2000

	ready int64
	agent exchange.Agent
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
	agent, status = exchange.NewDefaultAgent(PkgPath, messageHandler, false)
	if !status.OK() {
		fmt.Printf("init(\"%v\") failure: [%v]\n", PkgPath, status)
	}
	agent.Run()
}

var messageHandler core.MessageHandler = func(msg core.Message) {
	switch msg.Event {
	case core.StartupEvent:
		clientStartup(msg)
	case core.ShutdownEvent:
		clientShutdown()
	case core.PingEvent:
		start := time.Now()
		core.SendReply(msg, ping(nil).SetDuration(time.Since(start)))
	}
}
