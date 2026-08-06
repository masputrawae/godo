package repo

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Repo struct {
	db *sql.DB
	sq sq.StatementBuilderType
}

func New(db *sql.DB) *Repo {
	return &Repo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
	}
}

