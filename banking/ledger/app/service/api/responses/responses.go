package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

type Error struct {
	Error string `json:"error"`
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, Error{Error: msg})
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, Error{Error: "server error"})
}
