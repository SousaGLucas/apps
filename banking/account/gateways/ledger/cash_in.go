package ledger

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type CashInRequest struct {
	Amount int `json:"amount"`
}

func (c *Client) CashIn(ctx context.Context, ledgerAccountID uuid.UUID, amount int) error {
	const path = "/ledger/api/v1/accounts/%s/cash_in"

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(path, ledgerAccountID), CashInRequest{
		Amount: amount,
	})
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	_, err = c.executeRequest(req)
	if err != nil {
		var errResp ErrorResponse
		if errors.As(err, &errResp) {
			if errResp.Err == "account not found" {
				return entities.ErrAccountNotFound
			}
		}
		return fmt.Errorf("executing request: %w", err)
	}

	return nil
}
