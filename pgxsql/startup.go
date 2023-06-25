package pgxsql

import (
	"github.com/go-ai-agent/core/controller"
	"github.com/go-ai-agent/core/resource"
	"github.com/go-ai-agent/core/runtime"
	"reflect"
	"sync/atomic"
	"time"
)

type pkg struct{}

var (
	Uri     = pkgPath
	c       = make(chan resource.Message, 1)
	pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()
	started int64
	zone    = "zone"
	region  = "region"
)

// IsStarted - returns status of startup
func IsStarted() bool {
	return atomic.LoadInt64(&started) != 0
}

func setStarted() {
	atomic.StoreInt64(&started, 1)
}

func resetStarted() {
	atomic.StoreInt64(&started, 0)
}

func init() {
	resource.Register(Uri, c)
	go receive()
}

var messageHandler resource.MessageHandler = func(msg resource.Message) {
	switch msg.Event {
	case resource.StartupEvent:
		clientStartup(msg)
		//if IsStarted() {
		//	apply := resource.AccessControllerApply(&msg)
		//	if apply != nil {
		//		controllerApply = apply
		//	}
		//}
	case resource.ShutdownEvent:
		ClientShutdown()
	case resource.PingEvent:
		start := time.Now()
		resource.ReplyTo(msg, Ping[runtime.LogError, controller.DefaultHandler](nil).SetDuration(time.Since(start)))
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
