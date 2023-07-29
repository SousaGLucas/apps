package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/auth/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
	"github.com/SousaGLucas/apps/banking/auth/domain/usecases/user"
	"github.com/SousaGLucas/apps/banking/auth/telemetry"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func CreateUserHandler(uc user.CreateUserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		var body CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Error("parsing request body", zap.Error(err))

			responses.BadRequest(w, r, "invalid body")
			return
		}

		output, err := uc.CreateUser(ctx, user.CreateUserUseCaseInput{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			logger.Error("creating user", zap.Error(err))

			switch {
			case errors.Is(err, entities.ErrUserAlreadyExists):
				responses.BadRequest(w, r, "user already exists")
				return
			}

			responses.InternalServerError(w, r)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, CreateUserResponse{
			UserID: output.ID.String(),
		})
	}
}
