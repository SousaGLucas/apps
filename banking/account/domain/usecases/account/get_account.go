package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type GetAccountUseCase struct {
	DB GetAccountDB
}

type GetAccountInput struct {
	AccountID uuid.UUID
}

type GetAccountOutput struct {
	Account entities.Account
}

type GetAccountDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
}

func (uc GetAccountUseCase) GetAccount(ctx context.Context, input GetAccountInput) (GetAccountOutput, error) {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return GetAccountOutput{}, fmt.Errorf("retrieving account: %w", err)
	}

	return GetAccountOutput{
		Account: account,
	}, nil
}
