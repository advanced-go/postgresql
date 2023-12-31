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

[PostgresSQL][pgxsqlpkg] provides functions for query, exec, ping, and stat. Testing proxies are implemented for exec and query functions.
The processing of host generated messaging for startup and ping events is also supported. 

~~~
// Exec - templated function for executing a SQL statement
func Exec(ctx context.Context, req Request) (tag CommandTag, status runtime.Status) {
    // implementation details
}

// Query - function for a Query
func Query(ctx context.Context, req Request) (result pgx.Rows, status runtime.Status) {
// implementation details
}

// Ping - function for pinging the database cluster
func Ping(ctx context.Context) (status runtime.Status) {
// implementation details
}

// Stat - templated function for retrieving runtime stats
func Stat(ctx context.Context) (stat *pgxpool.Stat, status runtime.Status) {
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

Scanning of PostgreSQL rows into application types utilizes a templated interface, and corresponding templated Scan function.

~~~
// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows pgx.Rows) ([]T, error) {
    // implementation details
}
~~~

Resiliency and access logging for PostgresSQL database client calls is provided by an agent and controller types.


[rgriesemer]: <https://www.youtube.com/watch?v=0ReKdcpNyQg>
[pgxdmlpkg]: <https://pkg.go.dev/github.com/gotemplates/postgresql/pgxdml>
[pgxsqlpkg]: <https://pkg.go.dev/github.com/gotemplates/postgresql/pgxsql>
