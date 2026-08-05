package repo

import (
	"database/sql"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type Repo struct {
	db      *sql.DB
	sq      sq.StatementBuilderType
	table   string
	columns []string
}

func New(db *sql.DB) *Repo {
	slog.Info("🚩 init repo")

	return &Repo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question).RunWith(db),
	}
}
