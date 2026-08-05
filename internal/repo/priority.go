package repo

type Priority struct {
	*Repo
}

func (r *Repo) NewPriority() *Priority {
	repo := &Repo{
		table:   "statuses",
		columns: []string{"id", "emoji", "name"},
	}
	return &Priority{repo}
}
