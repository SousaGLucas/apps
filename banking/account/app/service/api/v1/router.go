package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	CreateAccountHandler http.HandlerFunc
	GetAccountHandler    http.HandlerFunc
	GetBalanceHandler    http.HandlerFunc
	DepositHandler       http.HandlerFunc
	WithdrawHandler      http.HandlerFunc
}

func (r *API) Routes(router *chi.Mux) {
	router.Post("/account/api/v1/accounts", r.CreateAccountHandler)
	router.Get("/account/api/v1/accounts/{account_id}", r.GetAccountHandler)
	router.Get("/account/api/v1/accounts/{account_id}/balance", r.GetBalanceHandler)
	router.Post("/account/api/v1/accounts/{account_id}/deposit", r.DepositHandler)
	router.Post("/account/api/v1/accounts/{account_id}/withdraw", r.WithdrawHandler)
}
