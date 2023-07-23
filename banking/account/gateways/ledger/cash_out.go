package ledger

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type CashOutRequest struct {
	Amount int `json:"amount"`
}

func (c *Client) CashOut(ctx context.Context, ledgerAccountID uuid.UUID, amount int) error {
	const path = "/api/v1/accounts/%s/cash_out"

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf(path, ledgerAccountID), CashOutRequest{
		Amount: amount,
	})
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	_, err = c.executeRequest(req)
	if err != nil {
		var errResp ErrorResponse
		if errors.As(err, &errResp) {
			switch errResp.Err {
			case "account not found":
				return entities.ErrAccountNotFound
			case "not enough funds":
				return entities.ErrNotEnoughFunds
			}
		}
		return fmt.Errorf("executing request: %w", err)
	}

	return nil
}
