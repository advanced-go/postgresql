package pgxsql

import (
	"github.com/go-ai-agent/core/runtime/startup"
	"time"
)

var (
	c = make(chan startup.Message, 1)
)

func init() {
	startup.Register(PkgUri, c)
	go receive()
}

var messageHandler startup.MessageHandler = func(msg startup.Message) {
	switch msg.Event {
	case startup.StartupEvent:
		clientStartup(msg)
	case startup.ShutdownEvent:
		ClientShutdown()
	case startup.PingEvent:
		start := time.Now()
		startup.ReplyTo(msg, Ping(nil).SetDuration(time.Since(start)))
	}
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
