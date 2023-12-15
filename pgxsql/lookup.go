package pgxsql

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
