package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type WithdrawUseCase struct {
	DB     WithdrawDB
	Ledger WithdrawLedgerClient
}

type WithdrawInput struct {
	AccountID uuid.UUID
	Amount    int
}

type WithdrawDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
}

type WithdrawLedgerClient interface {
	CashOut(ctx context.Context, ledgerAccountID uuid.UUID, amount int) error
}

func (uc WithdrawUseCase) Withdraw(ctx context.Context, input WithdrawInput) error {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return fmt.Errorf("retrieving account: %w", err)
	}

	err = uc.Ledger.CashOut(ctx, account.LedgerAccountID, input.Amount)
	if err != nil {
		return fmt.Errorf("creating ledger cash out: %w", err)
	}

	return nil
}
