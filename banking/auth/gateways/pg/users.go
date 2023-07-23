package pg

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SousaGLucas/apps/banking/auth/domain/entities"
	"github.com/SousaGLucas/apps/banking/auth/gateways/pg/sqlc"
)

type UsersRepository struct {
	DB *pgxpool.Pool
}

func (r UsersRepository) CreateUser(ctx context.Context, user entities.User) error {
	err := sqlc.New().CreateUser(ctx, r.DB, sqlc.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: user.ID,
			Valid: true,
		},
		Name:         user.Name,
		Email:        user.Email,
		Hashpassword: user.HashPassword,
		CreatedAt: pgtype.Timestamptz{
			Time:  user.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_email_key" {
				return entities.ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r UsersRepository) ListUsers(ctx context.Context) ([]entities.User, error) {
	rawUsers, err := sqlc.New().ListUsers(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	return sqlcListToUsers(rawUsers), nil
}

func (r UsersRepository) GetUser(ctx context.Context, userID uuid.UUID) (entities.User, error) {
	rawUser, err := sqlc.New().GetUsers(ctx, r.DB, pgtype.UUID{
		Bytes: userID,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, entities.ErrUserNotFound
		}
		return entities.User{}, err
	}

	return sqlcToUser(rawUser), nil
}

func sqlcListToUsers(l []sqlc.User) []entities.User {
	users := make([]entities.User, 0, len(l))
	for _, u := range l {
		users = append(users, sqlcToUser(u))
	}

	return users
}

func sqlcToUser(u sqlc.User) entities.User {
	return entities.User{
		ID:           u.ID.Bytes,
		Name:         u.Name,
		Email:        u.Email,
		HashPassword: u.Hashpassword,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    &u.UpdatedAt.Time,
	}
}
