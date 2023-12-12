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
	queryController     = newQueryController(queryControllerName, thresholdValues{}, nil)
	execControllerName  = "exec"
	execController      = newExecController(execControllerName, thresholdValues{}, nil)
	pingControllerName  = "ping"
	pingController      = newPingController(pingControllerName, thresholdValues{}, nil)
	statAgent           statusAgent
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
		if statAgent != nil {
			statAgent.Stop()
		}
		clientShutdown()
	case core.PingEvent:
		start := time.Now()
		core.SendReply(msg, ping(nil).SetDuration(time.Since(start)))
	}
}

func configControllers(msg core.Message) bool {
	// Need to also configure all controllers, query, exec and ping
	var err error
	statAgent, err = newStatusAgent(10, time.Second*2, &queryController, &execController)
	if err != nil {
		//Send error message
		return false
	}
	statAgent.Run()
	return true
}

// Scrap
//controllerApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
//	return func() {}, ctx, false
//}
