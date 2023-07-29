package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Account struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	LedgerAccountID uuid.UUID
	CreatedAt       time.Time
}

func NewAccount(userID uuid.UUID, ledgerAccountID uuid.UUID) (Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Account{}, fmt.Errorf("crating account id: %w", err)
	}

	return Account{
		ID:              id,
		UserID:          userID,
		LedgerAccountID: ledgerAccountID,
		CreatedAt:       time.Now(),
	}, nil
}
