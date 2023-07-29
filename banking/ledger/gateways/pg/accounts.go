package pg

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/gateways/pg/sqlc"
)

type AccountsRepository struct {
	DB *pgxpool.Pool
}

func (r AccountsRepository) CreateAccount(ctx context.Context, account entities.Account) error {
	err := sqlc.New().CreateAccount(ctx, r.DB, sqlc.CreateAccountParams{
		ID: pgtype.UUID{
			Bytes: account.ID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  account.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r AccountsRepository) GetAccount(ctx context.Context, accountID uuid.UUID) (entities.Account, error) {
	rawAccount, err := sqlc.New().GetAccount(ctx, r.DB, pgtype.UUID{
		Bytes: accountID,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Account{}, entities.ErrAccountNotFound
		}
		return entities.Account{}, err
	}

	return sqlcToAccount(rawAccount), nil
}

func sqlcToAccount(u sqlc.Account) entities.Account {
	return entities.Account{
		ID:        u.ID.Bytes,
		CreatedAt: u.CreatedAt.Time,
	}
}
