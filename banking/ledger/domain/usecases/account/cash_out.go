package account

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/domain/vos"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
)

type CashOutUseCase struct {
	DB       CashOutDB
	Balancer GetBalanceUseCase
}

type CashOutInput struct {
	AccountID uuid.UUID
	Amount    int
}

type CashOutDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
	CreateEvent(ctx context.Context, event entities.AccountEvent) error
}

func (uc CashOutUseCase) CashOut(ctx context.Context, input CashInInput) error {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return fmt.Errorf("retrieving account: %w", err)
	}

	balance, err := uc.Balancer.GetBalance(ctx, GetBalanceInput{
		AccountID: account.ID,
	})
	if err != nil {
		return fmt.Errorf("retrieving account balance: %w", err)
	}

	if balance.Balance < input.Amount {
		return ErrNotEnoughFunds
	}

	event, err := entities.NewAccountEvent(account, vos.AccountEventTypeCashOut, input.Amount)
	if err != nil {
		return fmt.Errorf("creating cash in event: %w", err)
	}

	err = uc.DB.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("saving cash in event: %w", err)
	}

	return nil
}
