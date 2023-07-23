package account

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/ledger/domain"
	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/domain/vos"
)

type GetBalanceUseCase struct {
	DB GetBalanceDB
}

type GetBalanceInput struct {
	AccountID uuid.UUID
}

type GetBalanceOutput struct {
	Balance int
}

type GetBalanceDB interface {
	GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error)
	ListAccountEvents(ctx context.Context, filter domain.ListAccountEventsFilter) ([]entities.AccountEvent, error)
}

func (uc GetBalanceUseCase) GetBalance(ctx context.Context, input GetBalanceInput) (GetBalanceOutput, error) {
	account, err := uc.DB.GetAccount(ctx, input.AccountID)
	if err != nil {
		return GetBalanceOutput{}, fmt.Errorf("retrieving account: %w", err)
	}

	const pageSize = 100

	var (
		balance       int
		lastFetchedID uuid.UUID
	)

	running := true
	for running {
		events, err := uc.DB.ListAccountEvents(ctx, domain.ListAccountEventsFilter{
			AccountID:     account.ID,
			LastFetchedID: lastFetchedID,
			PageSize:      pageSize,
		})
		if err != nil {
			return GetBalanceOutput{}, fmt.Errorf("listing account events: %w", err)
		}

		for _, event := range events {
			switch event.Type {
			case vos.AccountEventTypeCashIn:
				balance += event.Amount
			case vos.AccountEventTypeCashOut:
				balance -= event.Amount
			default:
				return GetBalanceOutput{}, fmt.Errorf("invalid event type: %w", err)
			}
		}

		if len(events) > 0 {
			lastFetchedID = events[len(events)-1].ID
		}

		if len(events) < pageSize {
			running = false
		}
	}

	return GetBalanceOutput{
		Balance: balance,
	}, nil
}
