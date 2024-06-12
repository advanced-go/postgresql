package pgxsql

import (
	"fmt"
	"github.com/advanced-go/stdlib/core"
	"net/url"
	"strconv"
)

func originFilter(values url.Values) []Entry {
	if values == nil {
		return storage
	}
	var result []Entry

	filter := core.NewOrigin(values)
	for _, e := range storage {
		target := core.Origin{Region: e.Region, Zone: e.Zone, SubZone: e.SubZone, Host: e.Host}
		if core.OriginMatch(target, filter) {
			result = append(result, e)
		}
	}
	if len(result) == 0 {
		return result
	}
	result = order(values, result)
	result = top(values, result)
	return distinct(values, result)
}

func order(values url.Values, entries []Entry) []Entry {
	if entries == nil || values == nil {
		return entries
	}
	s := values.Get("order")
	if s != "desc" {
		return entries
	}
	var result []Entry

	for index := len(entries) - 1; index >= 0; index-- {
		result = append(result, entries[index])
	}
	return result
}

func top(values url.Values, entries []Entry) []Entry {
	if entries == nil || values == nil {
		return entries
	}
	s := values.Get("top")
	if s == "" {
		return entries
	}
	cnt, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("top value is not valid: %v", s)
	}
	var result []Entry
	for i, e := range entries {
		if i < cnt {
			result = append(result, e)
		} else {
			break
		}
	}
	return result
}

func distinct(values url.Values, entries []Entry) []Entry {
	if entries == nil || values == nil {
		return entries
	}
	s := values.Get("distinct")
	if s == "" {
		return entries
	}
	if s != "host" {
		return entries
	}
	m := make(map[string]string)
	var result []Entry

	for _, e := range entries {
		_, ok := m[e.Host]
		if ok {
			continue
		}
		result = append(result, e)
		m[e.Host] = e.Host
	}
	return result
}
