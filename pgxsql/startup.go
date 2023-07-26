package pgxsql

import (
	"github.com/go-ai-agent/core/host"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/resiliency/controller"
	"reflect"
	"sync/atomic"
	"time"
)

type pkg struct{}

var (
	Uri     = pkgPath
	c       = make(chan host.Message, 1)
	pkgPath = reflect.TypeOf(any(pkg{})).PkgPath()
	started int64
	origin  = runtime.Origin{
		Region:     "region",
		Zone:       "zone",
		SubZone:    "",
		Service:    "",
		InstanceId: "",
	}
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
	host.Register(Uri, c)
	go receive()
}

var messageHandler host.MessageHandler = func(msg host.Message) {
	switch msg.Event {
	case host.StartupEvent:
		clientStartup(msg)
		//if IsStarted() {
		//	apply := host.AccessControllerApply(&msg)
		//	if apply != nil {
		//		controllerApply = apply
		//	}
		//}
	case host.ShutdownEvent:
		ClientShutdown()
	case host.PingEvent:
		start := time.Now()
		host.ReplyTo(msg, Ping[runtime.LogError, controller.DefaultHandler](nil).SetDuration(time.Since(start)))
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
