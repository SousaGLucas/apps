package entities

import (
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}
