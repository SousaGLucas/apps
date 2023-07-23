package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/SousaGLucas/apps/banking/account/domain/entities"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (c *Client) CreateUser(ctx context.Context, user entities.User) (uuid.UUID, error) {
	const path = "/api/v1/users"

	req, err := c.newRequest(ctx, http.MethodPost, path, CreateUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		var errResp ErrorResponse
		if errors.As(err, &errResp) {
			if errResp.Err == "user already exists" {
				return uuid.UUID{}, entities.ErrUserAlreadyExists
			}
		}
		return uuid.UUID{}, fmt.Errorf("executing request: %w", err)
	}

	var body CreateUserResponse
	if err := json.Unmarshal(resp, &body); err != nil {
		return uuid.UUID{}, fmt.Errorf("unmarshaling response body: %w", err)
	}

	userID, err := uuid.FromString(body.UserID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("parsing user id: %w", err)
	}

	return userID, nil
}
