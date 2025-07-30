package directus

import (
	"context"
	"encoding/json"
)

type AuthLogoutOptions struct {
	RefreshToken string   `json:"refresh_token"`
	Mode         AuthMode `json:"mode,omitempty"`
}

func (c *Client) AuthLogout(ctx context.Context, options AuthLogoutOptions) error {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Context-Type", "application/json").
		SetBody(options).
		Post("/auth/logout")

	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		return nil
	}

	var payload ErrorResponse
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		return err
	}

	return Error{
		Status:  resp.StatusCode(),
		Details: payload.Errors,
	}
}
