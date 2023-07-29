package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/ledger/domain/vos"
)

var (
	ErrInvalidEventAmount = errors.New("invalid event amount")
)

type AccountEvent struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Type      vos.AccountEventType
	Amount    int
	CreatedAt time.Time
}

func NewAccountEvent(acc Account, t vos.AccountEventType, amount int) (AccountEvent, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return AccountEvent{}, fmt.Errorf("crating account id: %w", err)
	}

	if amount <= 0 {
		return AccountEvent{}, ErrInvalidEventAmount
	}

	return AccountEvent{
		ID:        id,
		AccountID: acc.ID,
		Type:      t,
		Amount:    amount,
		CreatedAt: time.Now(),
	}, nil
}
