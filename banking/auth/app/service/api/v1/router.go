package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	Ping http.HandlerFunc

	CreateUserHandler http.HandlerFunc
	GetUserHandler    http.HandlerFunc
	ListUsersHandler  http.HandlerFunc
}

func (r *API) Routes(router *chi.Mux) {
	router.Get("/", r.Ping)

	router.Post("/api/v1/users", r.CreateUserHandler)
	router.Get("/api/v1/users", r.ListUsersHandler)
	router.Get("/api/v1/users/{id}", r.GetUserHandler)
}
