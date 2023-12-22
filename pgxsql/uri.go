package pgxsql

import (
	"github.com/advanced-go/core/uri"
)

var (
	lookup uri.Lookup
)

func init() {
	lookup = uri.NewLookup()
}
