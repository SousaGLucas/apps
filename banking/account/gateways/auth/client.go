package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	BaseURL string
}

type Client struct {
	baseURL *url.URL
	cli     *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base url: %w", err)
	}

	cli := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{
		baseURL: u,
		cli:     cli,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, data any) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing url path: %w", err)
	}

	var body bytes.Buffer
	if data != nil {
		if err = json.NewEncoder(&body).Encode(data); err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-type", "application/json")

	return req, nil
}

func (c *Client) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := c.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading http response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	default:
		var errResp ErrorResponse
		if err = json.Unmarshal(body, &errResp); err == nil {
			return nil, errResp
		}
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

type ErrorResponse struct {
	Err string `json:"error"`
}

func (e ErrorResponse) Error() string {
	return e.Err
}
