package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	Ping http.HandlerFunc

	CreateAccountHandler http.HandlerFunc
	GetAccountHandler    http.HandlerFunc
	GetBalanceHandler    http.HandlerFunc
	CashInHandler        http.HandlerFunc
	CashOutHandler       http.HandlerFunc
}

func (r *API) Routes(router *chi.Mux) {
	router.Get("/", r.Ping)

	router.Post("/api/v1/accounts", r.CreateAccountHandler)
	router.Get("/api/v1/accounts/{account_id}", r.GetAccountHandler)
	router.Get("/api/v1/accounts/{account_id}/balance", r.GetBalanceHandler)
	router.Post("/api/v1/accounts/{account_id}/cash_in", r.CashInHandler)
	router.Post("/api/v1/accounts/{account_id}/cash_out", r.CashOutHandler)
}
