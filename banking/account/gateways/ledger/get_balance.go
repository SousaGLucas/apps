package ledger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

func (c *Client) GetBalance(ctx context.Context, ledgerAccountID uuid.UUID) (int, error) {
	const path = "/ledger/api/v1/accounts/%s/balance"

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf(path, ledgerAccountID), nil)
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		var errResp ErrorResponse
		if errors.As(err, &errResp) {
			if errResp.Err == "account not found" {
				return 0, entities.ErrAccountNotFound
			}
		}
		return 0, fmt.Errorf("executing request: %w", err)
	}

	var body GetBalanceResponse
	if err := json.Unmarshal(resp, &body); err != nil {
		return 0, fmt.Errorf("unmarshaling response body: %w", err)
	}

	return body.Balance, nil
}
