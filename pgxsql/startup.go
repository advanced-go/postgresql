package pgxsql

import (
	"github.com/advanced-go/messaging/core"
	"time"
)

var (
	c                   = make(chan core.Message, 1)
	queryControllerName = "query"
	queryController     = NewQueryController(queryControllerName, Threshold{}, nil)
	execControllerName  = "exec"
	execController      = NewExecController(execControllerName, Threshold{}, nil)
	pingControllerName  = "exec"
	pingController      = NewPingController(pingControllerName, Threshold{}, nil)
	statusAgent         StatusAgent
)

func init() {
	core.Register(PkgUri, c)
	go receive()
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
		core.ReplyTo(msg, Ping(nil).SetDuration(time.Since(start)))
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

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			go messageHandler(msg)
		default:
		}
	}
}

// Scrap
//controllerApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
//	return func() {}, ctx, false
//}
