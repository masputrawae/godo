package repo

import (
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
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

func (r *Repo) IsSQLError(err error, code sqlite3.ErrNo) bool {
	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code == code
	}

	return false
}
