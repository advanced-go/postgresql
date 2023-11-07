package pgxsql

// Stats - pool statistics
//type Stats2 = pgxpool.Stat
/*
// CommandTag - results from an Exec command
type CommandTag struct {
	Sql          string
	RowsAffected int64
	Insert       bool
	Update       bool
	Delete       bool
	Select       bool
}

// FieldDescription - data for defining the returned Query fields/columns
type FieldDescription2 struct {
	Name                 string
	TableOID             uint32
	TableAttributeNumber uint16
	DataTypeOID          uint32
	DataTypeSize         int16
	TypeModifier         int32
	Format               int16
}

// Rows - interface that proxies the postgresql Rows interface
type Rows2 interface {
	// Close closes the rows, making the connection ready for use again. It is safe
	// to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() error

	// CommandTag returns the command tag from this query. It is only available after Rows is closed.
	CommandTag() CommandTag2

	FieldDescriptions() []FieldDescription2

	// Next prepares the next row for reading. It returns true if there is another
	// row and false if no more rows are available. It automatically closes rows
	// when all rows are read.
	Next() bool

	// Scan reads the values from the current row into dest values positionally.
	// dest can include pointers to core pgxsql, values implementing the Scanner
	// interface, and nil. nil will skip the value entirely. It is an error to
	// call Scan without first calling Next() and checking that it returned true.
	Scan(dest ...any) error

	// Values returns the decoded row values. As with Scan(), it is an error to
	// call Values without first calling Next() and checking that it returned
	// true.
	Values() ([]any, error)

	// RawValues returns the unparsed bytes of the row values. The returned data is only valid until the next Next
	// call or the Rows is closed.
	RawValues() [][]byte

	// Conn returns the underlying *Conn on which the queryv1 was executed. This may return nil if Rows did not come from a
	// *Conn (e.g. if it was created by RowsFromResultReader)
	// TODO : determine use case
	//Conn() *Conn
}

type proxyRows struct {
	pgxRows pgx.Rows
	fd      []FieldDescription2
}

func (r *proxyRows) Close() {
	if r != nil {
		r.pgxRows.Close()
	}
}

func (r *proxyRows) Err() error {
	if r == nil {
		return nil
	}
	return r.pgxRows.Err()
}

func (r *proxyRows) CommandTag() CommandTag2 {
	if r == nil {
		return CommandTag2{}
	}
	t := r.pgxRows.CommandTag()
	return CommandTag2{RowsAffected: t.RowsAffected(), Sql: t.String()}
}

func (r *proxyRows) FieldDescriptions() []FieldDescription2 {
	if r == nil {
		return nil
	}
	return r.fd
}

func (r *proxyRows) Next() bool {
	if r == nil {
		return false
	}
	return r.pgxRows.Next()
}

func (r *proxyRows) Scan(dest ...any) error {
	if r == nil {
		return nil
	}
	return r.pgxRows.Scan(dest)
}

func (r *proxyRows) Values() ([]any, error) {
	if r == nil {
		return nil, nil
	}
	return r.pgxRows.Values()
}

func (r *proxyRows) RawValues() [][]byte {
	if r == nil {
		return nil
	}
	return r.pgxRows.RawValues()
}

func createFieldDescriptions(fields []pgconn.FieldDescription) []FieldDescription2 {
	var result []FieldDescription2
	for _, f := range fields {
		result = append(result, FieldDescription2{Name: f.Name,
			TableOID:             f.TableOID,
			TableAttributeNumber: f.TableAttributeNumber,
			DataTypeOID:          f.DataTypeOID,
			DataTypeSize:         f.DataTypeSize,
			TypeModifier:         f.TypeModifier,
			Format:               f.Format})
	}
	return result
}


*/
