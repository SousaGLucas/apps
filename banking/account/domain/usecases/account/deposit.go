package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type DepositUseCase struct {
	DB     DepositDB
	Ledger DepositLedgerClient
}

type DepositInput struct {
	AccountID uuid.UUID
	Amount    int
}

type DepositDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
}

type DepositLedgerClient interface {
	CashIn(ctx context.Context, ledgerAccountID uuid.UUID, amount int) error
}

func (uc DepositUseCase) Deposit(ctx context.Context, input DepositInput) error {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return fmt.Errorf("retrieving account: %w", err)
	}

	err = uc.Ledger.CashIn(ctx, account.LedgerAccountID, input.Amount)
	if err != nil {
		return fmt.Errorf("creating ledger cash in: %w", err)
	}

	return nil
}
