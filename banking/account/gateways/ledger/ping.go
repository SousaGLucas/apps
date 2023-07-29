package ledger

import (
	"context"
	"fmt"
	"net/http"
)

type PingResponse string

func (c *Client) Ping(ctx context.Context) (string, error) {
	const path = "/"

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}

	return string(resp), nil
}
