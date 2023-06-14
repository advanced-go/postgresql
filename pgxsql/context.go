package pgxsql

import "github.com/go-sre/core/sql"

func findExecProxy(proxies []any) func(*sql.Request) (CommandTag, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*sql.Request) (CommandTag, error)); ok {
			return fn
		}
	}
	return nil
}

func findQueryProxy(proxies []any) func(*sql.Request) (Rows, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*sql.Request) (Rows, error)); ok {
			return fn
		}
	}
	return nil
}
