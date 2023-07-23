package entities

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
)

type LedgerAccount struct {
	ID        uuid.UUID
	CreatedAt time.Time
}
