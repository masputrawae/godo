package repo

type Todo struct {
	*Repo
}

func (r *Repo) NewTodo() *Todo {
	repo := &Repo{
		table: "todos",
		columns: []string{
			"id", "task", "done",
			"due", "created_at", "updated_at",
			"user_id", "status_id", "priority_id",
		},
	}
	return &Todo{repo}
}
