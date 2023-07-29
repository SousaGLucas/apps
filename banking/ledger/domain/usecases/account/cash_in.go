package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/domain/vos"
)

type CashInUseCase struct {
	DB CashInDB
}

type CashInInput struct {
	AccountID uuid.UUID
	Amount    int
}

type CashInDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
	CreateEvent(ctx context.Context, event entities.AccountEvent) error
}

func (uc CashInUseCase) CashIn(ctx context.Context, input CashInInput) error {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return fmt.Errorf("retrieving account: %w", err)
	}

	event, err := entities.NewAccountEvent(account, vos.AccountEventTypeCashIn, input.Amount)
	if err != nil {
		return fmt.Errorf("creating cash in event: %w", err)
	}

	err = uc.DB.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("saving cash in event: %w", err)
	}

	return nil
}
