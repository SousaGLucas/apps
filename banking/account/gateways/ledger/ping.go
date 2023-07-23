package ledger

import (
	"context"
	"encoding/json"
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

	var body PingResponse
	if err := json.Unmarshal(resp, &body); err != nil {
		return "", fmt.Errorf("unmarshaling response body: %w", err)
	}

	return string(body), nil
}
