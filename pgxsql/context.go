package pgxsql

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func findExecProxy(proxies []any) func(Request) (pgconn.CommandTag, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(Request) (pgconn.CommandTag, error)); ok {
			return fn
		}
	}
	return nil
}

func findQueryProxy(proxies []any) func(Request) (pgx.Rows, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(Request) (pgx.Rows, error)); ok {
			return fn
		}
	}
	return nil
}
