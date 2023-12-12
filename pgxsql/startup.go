package pgxsql

import (
	"fmt"
	"github.com/advanced-go/messaging/core"
	"github.com/advanced-go/messaging/exchange"
	"sync/atomic"
	"time"
)

var (
	queryControllerName = "query"
	queryController     = NewQueryController(queryControllerName, Threshold{}, nil)
	execControllerName  = "exec"
	execController      = NewExecController(execControllerName, Threshold{}, nil)
	pingControllerName  = "ping"
	pingController      = NewPingController(pingControllerName, Threshold{}, nil)
	statusAgent         StatusAgent
	agent               exchange.Agent
	ready               int64
	pingThreshold       = 500
	queryThreshold      = 2000
	execThreshold       = 2000
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
	status := exchange.Register(exchange.NewMailbox(PkgPath, false))
	if status.OK() {
		agent, status = exchange.NewAgent(PkgPath, messageHandler, nil, nil)
	}
	if !status.OK() {
		fmt.Printf("init() failure: [%v]\n", PkgPath)
	}
	agent.Run()
}

var messageHandler core.MessageHandler = func(msg core.Message) {
	switch msg.Event {
	case core.StartupEvent:
		if configControllers(msg) {
			clientStartup(msg)
		}
	case core.ShutdownEvent:
		if statusAgent != nil {
			statusAgent.Stop()
		}
		ClientShutdown()
	case core.PingEvent:
		start := time.Now()
		core.SendReply(msg, ping(nil).SetDuration(time.Since(start)))
	}
}

func configControllers(msg core.Message) bool {
	// Need to also configure all controllers, query, exec and ping
	var err error
	statusAgent, err = NewStatusAgent(10, time.Second*2, &queryController, &execController)
	if err != nil {
		//Send error message
		return false
	}
	statusAgent.Run()
	return true
}

// Scrap
//controllerApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
//	return func() {}, ctx, false
//}
