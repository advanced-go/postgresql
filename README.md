# postgresql

## pgxdml

[PostgresDML][pgxdmlpkg] implements types that build SQL statements based on the configured attributes. Support is also available for selecting
PostgreSQL functions for timestamps and next values.

~~~
// ExpandSelect - given a template, expand the template to build a WHERE clause if configured
func ExpandSelect(template string, where []Attr) (string, error) {
}

// WriteInsert - build a SQL insert statement with a VALUES list
func WriteInsert(sql string, values [][]any) (string, error) {
}

// WriteUpdate - build a SQL update statement, including SET and WHERE clauses
func WriteUpdate(sql string, attrs []Attr, where []Attr) (string, error) {
}

// WriteDelete - build a SQL delete statement with a WHERE clause
func WriteDelete(sql string, where []Attr) (string, error) {
}
~~~

## pgxsql

[PostgresSQL][pgxsqlpkg] provides the templated functions for query, exec, ping, and stat. Testing proxies are implemented for exec and query functions.
The processing of host generated messaging for startup and ping events is also supported. 

~~~
// Exec - templated function for executing a SQL statement
func Exec[E runtime.ErrorHandler](ctx context.Context, expectedCount int64, req *Request, args ...any) (tag CommandTag, status *runtime.Status) {
    // implementation details
}

// Query - templated function for a Query
func Query[E runtime.ErrorHandler](ctx context.Context, req *Request, args ...any) (result Rows, status *runtime.Status) {
// implementation details
}

// Ping - templated function for pinging the database cluster
func Ping[E runtime.ErrorHandler](ctx context.Context) (status *runtime.Status) {
// implementation details
}

// Stat - templated function for retrieving runtime stats
func Stat[E runtime.ErrorHandler](ctx context.Context) (stat *Stats, status *runtime.Status) {
// implementation details
}
~~~

The Request type is created as follows:

~~~
// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(resource string) string {
	return buildUri(PostgresNID, QueryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(resource string) string {
	return buildUri(PostgresNID, InsertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(resource string) string {
	return buildUri(PostgresNID, UpdateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(resource string) string {
	return buildUri(PostgresNID, DeleteNSS, resource)
}

~~~

Scanning of PostgreSQL rows into application types utilizes a templated interface, and corresponding templated Scan function. Care was taken to not leak
any direct references to PostgresSQL specific packages.

~~~
// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows Rows) ([]T, error) {
    // implementation details
}
~~~

Resiliency for PostgresSQL database client calls is provided by an [Actuator][actuatorcall] function call that is initialized by the host on startup:
~~~
var actuatorApply messaging.ActuatorApply

fn, ctx, limited = actuatorApply(ctx, statusCode func() int, req.Uri, runtime.ContextRequestId(ctx), "GET")
defer fn()
~~~

[pgxdmlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxdml/http>
[pgxsqlpkg]: <https://pkg.go.dev/github.com/idiomatic-go/postgresql/pgxsql>
[actuatorcall]: <https://pkg.go.dev/github.com/idiomatic-go/resiliency/actuator#EgressApply>
