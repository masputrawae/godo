package repo

type Status struct {
	*Repo
}

func (r *Repo) NewStatus() *Status {
	repo := &Repo{
		table:   "statuses",
		columns: []string{"id", "emoji", "name"},
	}
	return &Status{repo}
}
