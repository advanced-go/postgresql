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
	return r.index != len(resultSet)
}

func (r *accessRows) Scan(dest ...any) error {
	return nil
}
func (r *accessRows) Values() ([]any, error) {
	return []any{
		resultSet[r.index].StartTime,
		resultSet[r.index].Duration,
		resultSet[r.index].Traffic,
		resultSet[r.index].CreatedTS,

		resultSet[r.index].Region,
		resultSet[r.index].Zone,
		resultSet[r.index].SubZone,
		resultSet[r.index].Host,
		resultSet[r.index].InstanceId,

		resultSet[r.index].RequestId,
		resultSet[r.index].RelatesTo,
		resultSet[r.index].Protocol,
		resultSet[r.index].Method,
		resultSet[r.index].From,
		resultSet[r.index].To,
		resultSet[r.index].Url,
		resultSet[r.index].Path,

		resultSet[r.index].StatusCode,
		resultSet[r.index].Encoding,
		resultSet[r.index].Bytes,

		resultSet[r.index].Route,
		resultSet[r.index].RouteTo,
		resultSet[r.index].Timeout,
		resultSet[r.index].RateLimit,
		resultSet[r.index].RateBurst,
		resultSet[r.index].ReasonCode,
	}, nil
}
func (r *accessRows) RawValues() [][]byte { return nil }

var resultSet []Entry

func accessQuery(ctx context.Context, sql string, req *request) (pgx.Rows, error) {
	resultSet = accessFilter(req.values2)
	rows := NewAccessRows()
	return rows, nil
}
