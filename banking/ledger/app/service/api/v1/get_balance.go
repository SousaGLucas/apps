package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/ledger/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/domain/usecases/account"
	"github.com/SousaGLucas/apps/banking/ledger/telemetry"
)

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

func GetBalanceHandler(uc account.GetBalanceUseCase) http.HandlerFunc {
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

		output, err := uc.GetBalance(ctx, account.GetBalanceInput{
			AccountID: accountID,
		})
		if err != nil {
			logger.Error("retrieving balance", zap.Error(err))

			switch {
			case errors.Is(err, entities.ErrAccountNotFound):
				responses.BadRequest(w, r, "account not found")
				return
			}

			responses.InternalServerError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetBalanceResponse{
			Balance: output.Balance,
		})
	}
}
