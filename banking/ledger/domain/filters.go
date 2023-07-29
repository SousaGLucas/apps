package domain

import "github.com/gofrs/uuid/v5"

type ListAccountEventsFilter struct {
	AccountID     uuid.UUID
	LastFetchedID uuid.UUID
	PageSize      int
}
