package pgxsql

import (
	"github.com/advanced-go/stdlib/json"
	"github.com/jackc/pgx/v5/pgconn"
)

// CommandTag - results from an Exec command
type CommandTag struct {
	Sql          string `json:"sql"`
	RowsAffected int64  `json:"rows-affected"`
	Insert       bool   `json:"insert"`
	Update       bool   `json:"update"`
	Delete       bool   `json:"delete"`
	Select       bool   `json:"select"`
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

func NewCommandTag(url string) CommandTag {
	tag, status := json.New[CommandTag](url, nil)
	if !status.OK() {
		return CommandTag{}
	}
	return tag
}
