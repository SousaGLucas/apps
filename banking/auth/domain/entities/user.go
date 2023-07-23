package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	HashPassword string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func NewUser(name, email, hashPassword string) (User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, fmt.Errorf("crating user id: %w", err)
	}

	return User{
		ID:           id,
		Name:         name,
		Email:        email,
		HashPassword: hashPassword,
		CreatedAt:    time.Now(),
	}, nil
}
