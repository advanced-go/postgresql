package pgxsql

import "github.com/jackc/pgx/v5/pgconn"

// CommandTag - results from an Exec command
type CommandTag struct {
	Sql          string
	RowsAffected int64
	Insert       bool
	Update       bool
	Delete       bool
	Select       bool
}

func newCmdTag(tag pgconn.CommandTag) CommandTag {
	return CommandTag{
		Sql:          tag.String(),
		RowsAffected: tag.RowsAffected(),
		Insert:       tag.Insert(),
		Update:       tag.Update(),
		Delete:       tag.Delete(),
		Select:       tag.Select(),
	}
}
