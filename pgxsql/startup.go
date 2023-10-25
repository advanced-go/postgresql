package pgxsql

import (
	"context"
	"github.com/go-ai-agent/core/host"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

var (
	c               = make(chan host.Message, 1)
	controllerApply host.ControllerApply
)

func init() {
	controllerApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
		return func() {}, ctx, false
	}
	host.Register(PkgUri, c)
	go receive()
}

var messageHandler host.MessageHandler = func(msg host.Message) {
	switch msg.Event {
	case host.StartupEvent:
		clientStartup(msg)
		if IsStarted() {
			apply := host.AccessControllerApply(&msg)
			if apply != nil {
				controllerApply = apply
			}
		}
	case host.ShutdownEvent:
		ClientShutdown()
	case host.PingEvent:
		start := time.Now()
		host.ReplyTo(msg, Ping[runtime.LogError](nil).SetDuration(time.Since(start)))
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
