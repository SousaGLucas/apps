package user

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
)

type GetUserUseCase struct {
	DB GetUserDB
}

type GetUserInput struct {
	UserID uuid.UUID
}

type GetUserOutput struct {
	User entities.User
}

type GetUserDB interface {
	GetUser(ctx context.Context, userID uuid.UUID) (entities.User, error)
}

func (uc GetUserUseCase) GetUser(ctx context.Context, input GetUserInput) (GetUserOutput, error) {
	user, err := uc.DB.GetUser(ctx, input.UserID)
	if err != nil {
		return GetUserOutput{}, fmt.Errorf("retrieving user: %w", err)
	}

	return GetUserOutput{
		User: user,
	}, nil
}
