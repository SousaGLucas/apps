package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "Banking Auth Microservice V1")
	}
}
