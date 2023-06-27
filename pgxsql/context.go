package pgxsql

func findExecProxy(proxies []any) func(*Request) (CommandTag, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*Request) (CommandTag, error)); ok {
			return fn
		}
	}
	return nil
}

func findQueryProxy(proxies []any) func(*Request) (Rows, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*Request) (Rows, error)); ok {
			return fn
		}
	}
	return nil
}
