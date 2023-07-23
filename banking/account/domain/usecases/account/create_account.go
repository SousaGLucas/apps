package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type CreateAccountUseCase struct {
	DB     CreateAccountUseCaseDB
	Auth   CreateAccountAuthClient
	Ledger CreateAccountLedgerClient
}

type CreateAccountUseCaseInput struct {
	Name     string
	Email    string
	Password string
}

type CreateAccountUseCaseOutput struct {
	ID uuid.UUID
}

type CreateAccountUseCaseDB interface {
	CreateAccount(ctx context.Context, account entities.Account) error
}

type CreateAccountAuthClient interface {
	CreateUser(ctx context.Context, user entities.User) (uuid.UUID, error)
}

type CreateAccountLedgerClient interface {
	CreateAccount(ctx context.Context) (uuid.UUID, error)
}

func (uc CreateAccountUseCase) CreateAccount(ctx context.Context, input CreateAccountUseCaseInput) (CreateAccountUseCaseOutput, error) {
	userID, err := uc.Auth.CreateUser(ctx, entities.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return CreateAccountUseCaseOutput{}, fmt.Errorf("creating user: %w", err)
	}

	ledgerAccountID, err := uc.Ledger.CreateAccount(ctx)
	if err != nil {
		return CreateAccountUseCaseOutput{}, fmt.Errorf("creating ledger account: %w", err)
	}

	account, err := entities.NewAccount(userID, ledgerAccountID)
	if err != nil {
		return CreateAccountUseCaseOutput{}, fmt.Errorf("creating account %w", err)
	}

	err = uc.DB.CreateAccount(ctx, account)
	if err != nil {
		return CreateAccountUseCaseOutput{}, fmt.Errorf("saving account: %w", err)
	}

	return CreateAccountUseCaseOutput{
		ID: account.ID,
	}, nil
}
