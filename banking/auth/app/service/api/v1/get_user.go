package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/auth/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
	"github.com/SousaGLucas/apps/banking/auth/domain/usecases/user"
	"github.com/SousaGLucas/apps/banking/auth/telemetry"
)

type GetUserResponse struct {
	User GetUserUserResponse `json:"user"`
}

type GetUserUserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

func GetUserHandler(uc user.GetUserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		id, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			logger.Error("parsing request user id", zap.Error(err))

			responses.BadRequest(w, r, "invalid user id")
			return
		}

		output, err := uc.GetUser(ctx, user.GetUserInput{
			UserID: id,
		})
		if err != nil {
			logger.Error("retrieving user", zap.Error(err))

			switch {
			case errors.Is(err, entities.ErrUserNotFound):
				responses.BadRequest(w, r, "user not found")
				return
			}

			responses.InternalServerError(w, r)
			return
		}

		u := GetUserUserResponse{
			ID:        output.User.ID.String(),
			Name:      output.User.Name,
			Email:     output.User.Email,
			CreatedAt: output.User.CreatedAt.Format(time.RFC3339),
		}

		if output.User.UpdatedAt != nil {
			s := output.User.UpdatedAt.Format(time.RFC3339)
			u.UpdatedAt = &s
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, u)
	}
}
