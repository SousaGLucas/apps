package user

import (
	"context"
	"fmt"

	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
)

type ListUsersUseCase struct {
	DB ListUsersDB
}

type ListUsersOutput struct {
	Users []entities.User
}

type ListUsersDB interface {
	ListUsers(ctx context.Context) ([]entities.User, error)
}

func (uc ListUsersUseCase) ListUsers(ctx context.Context) (ListUsersOutput, error) {
	users, err := uc.DB.ListUsers(ctx)
	if err != nil {
		return ListUsersOutput{}, fmt.Errorf("retrieving users: %w", err)
	}

	return ListUsersOutput{
		Users: users,
	}, nil
}
