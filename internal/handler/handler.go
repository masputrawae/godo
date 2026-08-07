package handler

import "godo/internal/service"

type Handler struct {
	svUser    *service.User
	svSession *service.Session
	svTodo    *service.Todo
}

func New(svUser *service.User, svSession *service.Session, svTodo *service.Todo) *Handler {
	return &Handler{svUser, svSession, svTodo}
}
