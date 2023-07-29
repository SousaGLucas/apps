package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/auth/app/service/api/responses"
	"github.com/SousaGLucas/apps/banking/auth/domain/usecases/user"
	"github.com/SousaGLucas/apps/banking/auth/telemetry"
)

type ListUsersResponse struct {
	Users []ListUsersUsersResponse `json:"users"`
}

type ListUsersUsersResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

func ListUsersHandler(uc user.ListUsersUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger, err := telemetry.Logger(ctx)
		if err != nil {
			responses.InternalServerError(w, r)
			return
		}

		output, err := uc.ListUsers(ctx)
		if err != nil {
			logger.Error("listing users", zap.Error(err))

			responses.InternalServerError(w, r)
			return
		}

		users := make([]ListUsersUsersResponse, 0, len(output.Users))

		for _, item := range output.Users {
			u := ListUsersUsersResponse{
				ID:        item.ID.String(),
				Name:      item.Name,
				Email:     item.Email,
				CreatedAt: item.CreatedAt.Format(time.RFC3339),
				UpdatedAt: nil,
			}

			if u.UpdatedAt != nil {
				s := item.UpdatedAt.Format(time.RFC3339)
				u.UpdatedAt = &s
			}

			users = append(users, u)
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, ListUsersResponse{
			Users: users,
		})
	}
}
