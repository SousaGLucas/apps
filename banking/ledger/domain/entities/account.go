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
	ID        uuid.UUID
	CreatedAt time.Time
}

func NewAccount() (Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Account{}, fmt.Errorf("crating account id: %w", err)
	}

	return Account{
		ID:        id,
		CreatedAt: time.Now(),
	}, nil
}
