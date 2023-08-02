package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	CreateUserHandler http.HandlerFunc
	GetUserHandler    http.HandlerFunc
	ListUsersHandler  http.HandlerFunc
}

func (r *API) Routes(router *chi.Mux) {
	router.Post("/auth/api/v1/users", r.CreateUserHandler)
	router.Get("/auth/api/v1/users", r.ListUsersHandler)
	router.Get("/auth/api/v1/users/{id}", r.GetUserHandler)
}
