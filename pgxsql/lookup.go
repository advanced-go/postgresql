package pgxsql

import (
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
)

var (
	overrideLookup func(string) []string
)

func setOverrideLookup(t []string) {
	if t == nil {
		overrideLookup = nil
		return
	}
	overrideLookup = func(key string) []string { return t }
}

func lookup(key string) ([]string, bool) {
	if overrideLookup == nil || len(key) == 0 {
		return nil, false
	}
	val := overrideLookup(key)
	if len(val) > 0 {
		return val, true
	}
	return nil, false
}

func execOverride(urls []string) (CommandTag, runtime.Status) {
	if len(urls) == 0 {
		return CommandTag{}, runtime.NewStatus(runtime.StatusInvalidArgument)
	}
	if len(urls) == 1 {
		return io2.ReadState[CommandTag](urls[0])
	}
	tag := CommandTag{}
	status := runtime.StatusOK()
	if len(urls[0]) > 0 {
		tag, status = io2.ReadState[CommandTag](urls[0])
		if !status.OK() {
			return CommandTag{}, status
		}
	}
	return tag, io2.ReadStatus(urls[1])
}
