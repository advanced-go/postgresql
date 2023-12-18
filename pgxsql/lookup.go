package pgxsql

import "github.com/advanced-go/core/runtime"

var (
	overrideLookup func(string) []string
)

func setOverrideLookup(t any) {
	if t == nil {
		overrideLookup = nil
		return
	}
	overrideLookup = runtime.LookupFromType[func(string) []string](t)
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
