package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"time"
)

// PingFn - ping function
type PingFn func(ctx context.Context) *runtime.Status

// StatusAgent - an agent that will manage returning an endpoint back to receiving traffic
type StatusAgent interface {
	Run()
	Stop()
}

type agentConfig struct {
	interval  time.Duration
	reset     time.Duration
	threshold int
	quit      chan struct{}
	query     *QueryController
	exec      *ExecController
}

// NewStatusAgent - creation of an agent with configuration
func NewStatusAgent(threshold int, interval time.Duration, query *QueryController, exec *ExecController) (StatusAgent, error) {
	if interval <= 0 {
		return nil, errors.New(fmt.Sprintf("error: interval is less than or equal to zero [%v]", interval))
	}
	a := new(agentConfig)
	a.threshold = threshold
	a.interval = interval
	a.reset = time.Hour * 24
	a.quit = make(chan struct{}, 1)
	a.query = query
	a.exec = exec
	return a, nil
}

// Run - run the agent
func (a *agentConfig) Run() {
	//go run(a.threshold, a.interval, a.reset,a.quit, a.query, a.exec)
}

// Stop - stop the agent
func (a *agentConfig) Stop() {
	//a.quit <- struct{}{}
}

func run(threshold int, interval, r time.Duration, quit <-chan struct{}, query *QueryController, exec *ExecController) {
	tick := time.Tick(interval)
	reset := time.Tick(r)
	var ticks int64
	var failures int64
	var status runtime.Status

	for {
		select {
		case <-tick:
			ticks++
			status = ping(nil)
			if !status.OK() {
				failures++
			}
			if failures > int64(threshold) {
				// Need to decrease the rates on Query and Exec
			}
		case <-reset:
		default:
		}
		select {
		case <-quit:
			return
		default:
		}
	}
}
