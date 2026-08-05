package repo

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

type UserCreatePayload struct {
	Username string
	Password string
	Email    string
}
