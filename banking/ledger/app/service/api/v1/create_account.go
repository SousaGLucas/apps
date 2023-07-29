package v1

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/ledger/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/ledger/domain/usecases/account"
	"github.com/SousaGLucas/apps/banking/ledger/telemetry"
)

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}

func CreateAccountHandler(uc account.CreateAccountUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		output, err := uc.CreateAccount(ctx)
		if err != nil {
			logger.Error("creating account", zap.Error(err))

			responses.InternalServerError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, CreateAccountResponse{
			AccountID: output.ID.String(),
		})
	}
}
