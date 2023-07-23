// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamptz
}

type AccountEvent struct {
	ID        pgtype.UUID
	AccountID pgtype.UUID
	Type      string
	Amount    int32
	CreatedAt pgtype.Timestamptz
}