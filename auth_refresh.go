package directus

import (
	"context"
	"encoding/json"
)

type AuthRefreshOptions struct {
	RefreshToken string   `json:"refresh_token"`
	Mode         AuthMode `json:"mode,omitempty"`
}

func (c *Client) AuthRefresh(ctx context.Context, options AuthRefreshOptions) (AuthResult, error) {
	var payload AuthResponsePayload

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Context-Type", "application/json").
		SetBody(options).
		Post("/auth/refresh")

	if err != nil {
		return payload.Data, err
	}

	if resp.IsError() {
		var failed ErrorsPayload
		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
			return payload.Data, err
		}

		return payload.Data, failed.Errors
	}

	err = json.Unmarshal(resp.Body(), &payload)
	return payload.Data, err
}
