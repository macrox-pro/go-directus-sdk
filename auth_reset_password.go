package directus

import (
	"context"
	"encoding/json"
)

type PasswordResetRequestParams struct {
	Email string   `json:"email"`
	Mode  AuthMode `json:"mode,omitempty"`
}

type PasswordResetParams struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (c *Client) AuthResetPasswordRequest(ctx context.Context, options PasswordResetRequestParams) error {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(options).
		Post("/auth/password/request")

	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		return nil
	}

	var payload ErrorsPayload
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		return err
	}

	return payload.Errors
}

func (c *Client) AuthPasswordReset(ctx context.Context, options PasswordResetParams) error {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(options).
		Post("/auth/password/reset")

	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		return nil
	}

	var payload ErrorsPayload
	if err := json.Unmarshal(resp.Body(), &payload); err != nil {
		return err
	}

	return payload.Errors
}
