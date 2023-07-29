package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
)

type CreateUserUseCase struct {
	DB CreateUserUseCaseDB
}

type CreateUserUseCaseInput struct {
	Name     string
	Email    string
	Password string
}

type CreateUserUseCaseOutput struct {
	ID uuid.UUID
}

type CreateUserUseCaseDB interface {
	CreateUser(ctx context.Context, user entities.User) error
}

func (uc CreateUserUseCase) CreateUser(ctx context.Context, input CreateUserUseCaseInput) (CreateUserUseCaseOutput, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(input.Password))
	if err != nil {
		return CreateUserUseCaseOutput{}, fmt.Errorf("creating hash password: %w", err)
	}

	user, err := entities.NewUser(input.Name, input.Email, fmt.Sprintf("%x", hash.Sum(nil)))
	if err != nil {
		return CreateUserUseCaseOutput{}, fmt.Errorf("creating user %w", err)
	}

	err = uc.DB.CreateUser(ctx, user)
	if err != nil {
		return CreateUserUseCaseOutput{}, fmt.Errorf("saving user: %w", err)
	}

	return CreateUserUseCaseOutput{
		ID: user.ID,
	}, nil
}
