package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type GetBalanceUseCase struct {
	DB     GetBalanceDB
	Ledger GetBalanceLedgerClient
}

type GetBalanceInput struct {
	AccountID uuid.UUID
}

type GetBalanceOutput struct {
	Balance int
}

type GetBalanceDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
}

type GetBalanceLedgerClient interface {
	GetBalance(ctx context.Context, ledgerAccountID uuid.UUID) (int, error)
}

func (uc GetBalanceUseCase) GetBalance(ctx context.Context, input GetBalanceInput) (GetBalanceOutput, error) {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return GetBalanceOutput{}, fmt.Errorf("retrieving account: %w", err)
	}

	balance, err := uc.Ledger.GetBalance(ctx, account.LedgerAccountID)
	if err != nil {
		return GetBalanceOutput{}, fmt.Errorf("retrieving account balance: %w", err)
	}

	return GetBalanceOutput{
		Balance: balance,
	}, nil
}
