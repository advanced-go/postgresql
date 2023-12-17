package pgxsql

var (
	overrideLookup func(string) []string
)

func setOverrideLookup(t any) {
	if t == nil {
		overrideLookup = nil
		return
	}
	if s, ok := t.(string); ok {
		overrideLookup = func(key string) []string { return []string{s, ""} }
		return
	}
	if s, ok := t.([]string); ok {
		overrideLookup = func(key string) []string { return s }
		return
	}
	if m, ok := t.(map[string][]string); ok {
		overrideLookup = func(key string) []string { return m[key] }
		return
	}
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
