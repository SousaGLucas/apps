package v1

import (
	"errors"
	"net/http"

	"encoding/json"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/account/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/account/domain/entities"
	"github.com/SousaGLucas/apps/banking/account/domain/usecases/account"
	"github.com/SousaGLucas/apps/banking/account/telemetry"
)

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

		var body CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Error("parsing request body", zap.Error(err))

			responses.BadRequest(w, r, "invalid body")
			return
		}

		output, err := uc.CreateAccount(ctx, account.CreateAccountUseCaseInput{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			logger.Error("creating account", zap.Error(err))

			switch {
			case errors.Is(err, entities.ErrAccountNotFound):
				responses.BadRequest(w, r, "account not found")
				return
			case errors.Is(err, entities.ErrUserAlreadyExists):
				responses.BadRequest(w, r, "user already exists")
				return
			}

			responses.InternalServerError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, CreateAccountResponse{
			AccountID: output.ID.String(),
		})
	}
}
