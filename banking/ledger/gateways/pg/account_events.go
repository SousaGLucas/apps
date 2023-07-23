package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SousaGLucas/apps/banking/ledger/domain"
	"github.com/SousaGLucas/apps/banking/ledger/domain/entities"
	"github.com/SousaGLucas/apps/banking/ledger/gateways/pg/sqlc"
)

type AccountEventsRepository struct {
	DB *pgxpool.Pool
}

func (r AccountEventsRepository) CreateEvent(ctx context.Context, event entities.AccountEvent) error {
	err := sqlc.New().CreateEvent(ctx, r.DB, sqlc.CreateEventParams{
		ID: pgtype.UUID{
			Bytes: event.ID,
			Valid: true,
		},
		AccountID: pgtype.UUID{
			Bytes: event.AccountID,
			Valid: true,
		},
		Type:   event.Type.String(),
		Amount: int32(event.Amount),
		CreatedAt: pgtype.Timestamptz{
			Time:  event.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r AccountEventsRepository) ListAccountEvents(ctx context.Context, filter domain.ListAccountEventsFilter) ([]entities.AccountEvent, error) {
	builder := psql.
		Select(
			"id",
			"account_id",
			"type",
			"amount",
			"created_at",
		).
		From("account_events").
		Where("account_id = ?", filter.AccountID.String()).
		OrderBy("id desc").
		Limit(uint64(filter.PageSize))

	if !filter.LastFetchedID.IsNil() {
		builder = builder.Where("id < ?", filter.LastFetchedID.String())
	}

	sql, params, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}

	events := make([]entities.AccountEvent, 0, filter.PageSize)

	for rows.Next() {
		var event entities.AccountEvent
		err := rows.Scan(
			&event.ID,
			&event.AccountID,
			&event.Type,
			&event.Amount,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
