package pgxsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/stdlib/access"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		case StartTimeName:
			log.StartTime = values[i].(time.Time)
		case DurationName:
			log.Duration = values[i].(int64)
		case TrafficName:
			log.Traffic = values[i].(string)
		case CreatedTSName:
			log.CreatedTS = values[i].(time.Time)

		case RegionName:
			log.Region = values[i].(string)
		case ZoneName:
			log.Zone = values[i].(string)
		case SubZoneName:
			log.SubZone = values[i].(string)
		case HostName:
			log.Host = values[i].(string)
		case InstanceIdName:
			log.InstanceId = values[i].(string)

		case RequestIdName:
			log.RequestId = values[i].(string)
		case RelatesToName:
			log.RelatesTo = values[i].(string)
		case ProtocolName:
			log.Protocol = values[i].(string)
		case MethodName:
			log.Method = values[i].(string)
		case FromName:
			log.From = values[i].(string)
		case ToName:
			log.To = values[i].(string)
		case UrlName:
			log.Url = values[i].(string)
		case PathName:
			log.Path = values[i].(string)

		case StatusCodeName:
			log.StatusCode = values[i].(int32)
		case EncodingName:
			log.Encoding = values[i].(string)
		case BytesName:
			log.Bytes = values[i].(int64)

		case RouteName:
			log.Route = values[i].(string)
		case RouteToName:
			log.RouteTo = values[i].(string)
		case TimeoutName:
			log.Timeout = values[i].(int32)
		case RateLimitName:
			log.RateLimit = values[i].(float64)
		case RateBurstName:
			log.RateBurst = values[i].(int32)
		case ReasonCodeName:
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

var list = []Entry{
	{time.Now().UTC(), 100, access.EgressTraffic, time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "www.google.com", "", "https://www.google.com/search?q-golang", "/search", 200, "gzip", 12345, "google-search", "primary", 500, 100, 10, ""},
	{time.Now().UTC(), 100, access.IngressTraffic, time.Now().UTC(), "us-west", "oregon", "dc1", "localhost:8081", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "github/advanced-go/search", "", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 504, "gzip", 12345, "search", "primary", 500, 100, 10, "TO"},
}

var fields = []pgconn.FieldDescription{
	{StartTimeName, 0, 0, 0, 0, 0, 0},
	{DurationName, 0, 0, 0, 0, 0, 0},
	{TrafficName, 0, 0, 0, 0, 0, 0},
	{CreatedTSName, 0, 0, 0, 0, 0, 0},

	{RegionName, 0, 0, 0, 0, 0, 0},
	{ZoneName, 0, 0, 0, 0, 0, 0},
	{SubZoneName, 0, 0, 0, 0, 0, 0},
	{HostName, 0, 0, 0, 0, 0, 0},
	{InstanceIdName, 0, 0, 0, 0, 0, 0},

	{RequestIdName, 0, 0, 0, 0, 0, 0},
	{RelatesToName, 0, 0, 0, 0, 0, 0},
	{ProtocolName, 0, 0, 0, 0, 0, 0},
	{MethodName, 0, 0, 0, 0, 0, 0},
	{FromName, 0, 0, 0, 0, 0, 0},
	{ToName, 0, 0, 0, 0, 0, 0},
	{UrlName, 0, 0, 0, 0, 0, 0},
	{PathName, 0, 0, 0, 0, 0, 0},

	{StatusCodeName, 0, 0, 0, 0, 0, 0},
	{EncodingName, 0, 0, 0, 0, 0, 0},
	{BytesName, 0, 0, 0, 0, 0, 0},

	{RouteName, 0, 0, 0, 0, 0, 0},
	{RouteToName, 0, 0, 0, 0, 0, 0},
	{TimeoutName, 0, 0, 0, 0, 0, 0},
	{RateLimitName, 0, 0, 0, 0, 0, 0},
	{RateBurstName, 0, 0, 0, 0, 0, 0},
	{ReasonCodeName, 0, 0, 0, 0, 0, 0},
}

type accessRows struct {
	index int
}

func NewAccessRows() pgx.Rows {
	a := new(accessRows)
	a.index = -1
	return a
}

func (r *accessRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r *accessRows) FieldDescriptions() []pgconn.FieldDescription {
	return fields
}

func (r *accessRows) Conn() *pgx.Conn {
	//TODO implement me
	return nil
}

func (r *accessRows) Close()     {}
func (r *accessRows) Err() error { return nil }

func (r *accessRows) Next() bool {
	r.index++
	return r.index != len(list)
}

func (r *accessRows) Scan(dest ...any) error {
	return nil
}
func (r *accessRows) Values() ([]any, error) {
	return []any{
		list[r.index].StartTime,
		list[r.index].Duration,
		list[r.index].Traffic,
		list[r.index].CreatedTS,

		list[r.index].Region,
		list[r.index].Zone,
		list[r.index].SubZone,
		list[r.index].Host,
		list[r.index].InstanceId,

		list[r.index].RequestId,
		list[r.index].RelatesTo,
		list[r.index].Protocol,
		list[r.index].Method,
		list[r.index].From,
		list[r.index].To,
		list[r.index].Url,
		list[r.index].Path,

		list[r.index].StatusCode,
		list[r.index].Encoding,
		list[r.index].Bytes,

		list[r.index].Route,
		list[r.index].RouteTo,
		list[r.index].Timeout,
		list[r.index].RateLimit,
		list[r.index].RateBurst,
		list[r.index].ReasonCode,
	}, nil
}
func (r *accessRows) RawValues() [][]byte { return nil }

func accessQuery(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows := NewAccessRows()
	return rows, nil
}
