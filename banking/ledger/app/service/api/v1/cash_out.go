package v1

import (
	"errors"
	"net/http"

	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/ledger/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/domain/usecases/account"
	"github.com/SousaGLucas/apps/banking/ledger/telemetry"
)

type CashOutRequest struct {
	Amount int `json:"amount"`
}

func CashOutHandler(uc account.CashOutUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		accountID, err := uuid.FromString(chi.URLParam(r, "account_id"))
		if err != nil {
			logger.Error("parsing request account id", zap.Error(err))

			responses.BadRequest(w, r, "invalid account id")
			return
		}

		var body CashInRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Error("parsing request body", zap.Error(err))

			responses.BadRequest(w, r, "invalid body")
			return
		}

		err = uc.CashOut(ctx, account.CashInInput{
			AccountID: accountID,
			Amount:    body.Amount,
		})
		if err != nil {
			logger.Error("creating cash in", zap.Error(err))

			switch {
			case errors.Is(err, entities.ErrAccountNotFound):
				responses.BadRequest(w, r, "account not found")
				return
			case errors.Is(err, account.ErrNotEnoughFunds):
				responses.BadRequest(w, r, "not enough funds")
				return
			}

			responses.InternalServerError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
	}
}
