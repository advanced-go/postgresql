package pgxsql

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
	"time"
)

var (
	c               = make(chan startup.Message, 1)
	controllerApply startup.ControllerApply
)

func init() {
	controllerApply = func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
		return func() {}, ctx, false
	}
	startup.Register(PkgUri, c)
	go receive()
}

var messageHandler startup.MessageHandler = func(msg startup.Message) {
	switch msg.Event {
	case startup.StartupEvent:
		clientStartup(msg)
		if IsStarted() {
			apply := startup.AccessControllerApply(&msg)
			if apply != nil {
				controllerApply = apply
			}
		}
	case startup.ShutdownEvent:
		ClientShutdown()
	case startup.PingEvent:
		start := time.Now()
		startup.ReplyTo(msg, Ping[runtime.LogError](nil).SetDuration(time.Since(start)))
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
