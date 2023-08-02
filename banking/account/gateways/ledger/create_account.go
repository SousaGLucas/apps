package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}

func (c *Client) CreateAccount(ctx context.Context) (uuid.UUID, error) {
	const path = "/ledger/api/v1/accounts"

	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("executing request: %w", err)
	}

	var body CreateAccountResponse
	if err := json.Unmarshal(resp, &body); err != nil {
		return uuid.UUID{}, fmt.Errorf("unmarshaling response body: %w", err)
	}

	accountID, err := uuid.FromString(body.AccountID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parsing account id: %w", err)
	}

	return accountID, nil
}
