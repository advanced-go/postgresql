package access

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/postgresql/pgxsql"
	"github.com/advanced-go/stdlib/access"
	"time"
)

// Entry - timeseries access log struct
type Entry struct {
	StartTime time.Time `json:"start-time"`
	Duration  int64     `json:"duration"`
	Traffic   string    `json:"traffic"`
	CreatedTS time.Time `json:"created-ts"`

	Region     string `json:"region"`
	Zone       string `json:"zone"`
	SubZone    string `json:"sub-zone"`
	Host       string `json:"host"`
	InstanceId string `json:"instance-id"`

	RequestId string `json:"request-id"`
	RelatesTo string `json:"relates-to"`
	Protocol  string `json:"proto"`
	Method    string `json:"method"`
	From      string `json:"from"`
	To        string `json:"to"`
	Url       string `json:"url"`
	Path      string `json:"path"`

	StatusCode int32  `json:"status-code"`
	Encoding   string `json:"encoding"`
	Bytes      int64  `json:"bytes"`

	Route      string  `json:"route"`
	RouteTo    string  `json:"route-to"`
	Timeout    int32   `json:"timeout"`
	RateLimit  float64 `json:"rate-limit"`
	RateBurst  int32   `json:"rate-burst"`
	ReasonCode string  `json:"rc"`
}

func (Entry) Scan(columnNames []string, values []any) (log Entry, err error) {
	for i, name := range columnNames {
		switch name {
		case pgxsql.StartTimeName:
			log.StartTime = values[i].(time.Time)
		case pgxsql.DurationName:
			log.Duration = values[i].(int64)
		case pgxsql.TrafficName:
			log.Traffic = values[i].(string)
		case pgxsql.CreatedTSName:
			log.CreatedTS = values[i].(time.Time)

		case pgxsql.RegionName:
			log.Region = values[i].(string)
		case pgxsql.ZoneName:
			log.Zone = values[i].(string)
		case pgxsql.SubZoneName:
			log.SubZone = values[i].(string)
		case pgxsql.HostName:
			log.Host = values[i].(string)
		case pgxsql.InstanceIdName:
			log.InstanceId = values[i].(string)

		case pgxsql.RequestIdName:
			log.RequestId = values[i].(string)
		case pgxsql.RelatesToName:
			log.RelatesTo = values[i].(string)
		case pgxsql.ProtocolName:
			log.Protocol = values[i].(string)
		case pgxsql.MethodName:
			log.Method = values[i].(string)
		case pgxsql.FromName:
			log.From = values[i].(string)
		case pgxsql.ToName:
			log.To = values[i].(string)
		case pgxsql.UrlName:
			log.Url = values[i].(string)
		case pgxsql.PathName:
			log.Path = values[i].(string)

		case pgxsql.StatusCodeName:
			log.StatusCode = values[i].(int32)
		case pgxsql.EncodingName:
			log.Encoding = values[i].(string)
		case pgxsql.BytesName:
			log.Bytes = values[i].(int64)

		case pgxsql.RouteName:
			log.Route = values[i].(string)
		case pgxsql.RouteToName:
			log.RouteTo = values[i].(string)
		case pgxsql.TimeoutName:
			log.Timeout = values[i].(int32)
		case pgxsql.RateLimitName:
			log.RateLimit = values[i].(float64)
		case pgxsql.RateBurstName:
			log.RateBurst = values[i].(int32)
		case pgxsql.ReasonCodeName:
			log.ReasonCode = values[i].(string)
		default:
			err = errors.New(fmt.Sprintf("invalid field name: %v", name))
			return
		}
	}
	return
}

func (a Entry) Values() []any {
	return []any{
		a.StartTime,
		a.Duration,
		a.Traffic,
		a.CreatedTS,

		a.Region,
		a.Zone,
		a.SubZone,
		a.Host,
		a.InstanceId,

		a.RequestId,
		a.RelatesTo,
		a.Protocol,
		a.Method,
		a.From,
		a.To,
		a.Url,
		a.Path,

		a.StatusCode,
		a.Encoding,
		a.Bytes,

		a.Route,
		a.RouteTo,
		a.Timeout,
		a.RateLimit,
		a.RateBurst,
		a.ReasonCode,
	}
}

func (Entry) Rows(entries []Entry) [][]any {
	var values [][]any

	for _, e := range entries {
		values = append(values, e.Values())
	}
	return values
}

var storage = []Entry{
	{time.Date(2024, 6, 10, 7, 10, 35, 0, time.UTC), 100, access.EgressTraffic, time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "www.google.com", "", "https://www.google.com/search?q-golang", "/search", 200, "gzip", 12345, "google-search", "primary", 500, 100, 10, ""},
	{time.Date(2024, 6, 10, 7, 12, 35, 0, time.UTC), 100, access.IngressTraffic, time.Now().UTC(), "us-west", "oregon", "dc2", "localhost:8081", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "github/advanced-go/search", "", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 504, "gzip", 12345, "search", "primary", 500, 100, 10, "TO"},
	{time.Date(2024, 6, 11, 8, 45, 2, 0, time.UTC), 100, access.EgressTraffic, time.Now().UTC(), "us-west", "oregon", "dc3", "www.test-host.com", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "www.google.com", "", "https://www.google.com/search?q-golang", "/search", 200, "gzip", 12345, "google-search", "primary", 500, 100, 10, ""},
	{time.Date(2024, 6, 12, 16, 55, 15, 0, time.UTC), 100, access.IngressTraffic, time.Now().UTC(), "us-west", "oregon", "dc4", "localhost:8080", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "github/advanced-go/search", "", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 504, "gzip", 12345, "search", "primary", 500, 100, 10, "TO"},
}

func accessInsert(ctx context.Context, sql string, req *pgxsql.request) (pgxsql.CommandTag, error) {
	if req == nil || req.values == nil {
		return pgxsql.CommandTag{}, errors.New("request or request values is nil")
	}
	var entry Entry
	var count = 0
	names := pgxsql.createColumnNames(fields)

	for _, row := range req.values {
		e, err := entry.Scan(names, row)
		if err == nil {
			count++
			storage = append(storage, e)
		}
	}
	return pgxsql.CommandTag{RowsAffected: int64(count), Insert: true}, nil
}
