package pgxsql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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
