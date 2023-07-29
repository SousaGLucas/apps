package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
)

type CreateAccountUseCase struct {
	DB CreateAccountUseCaseDB
}

type CreateAccountUseCaseOutput struct {
	ID uuid.UUID
}

type CreateAccountUseCaseDB interface {
	CreateAccount(ctx context.Context, account entities.Account) error
}

func (uc CreateAccountUseCase) CreateAccount(ctx context.Context) (CreateAccountUseCaseOutput, error) {
	account, err := entities.NewAccount()
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
