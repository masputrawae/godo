package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type repo struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func New(db *sql.DB) *repo {
	return &repo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
	}
}
