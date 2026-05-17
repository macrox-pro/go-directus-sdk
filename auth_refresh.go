package directus

import (
	"context"
	"encoding/json"
)

type AuthRefreshParams struct {
	RefreshToken string   `json:"refresh_token"`
	Mode         AuthMode `json:"mode,omitempty"`
}

func (c *Client) AuthRefresh(ctx context.Context, options AuthRefreshParams) (AuthResult, error) {
	var payload AuthResponsePayload

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader(headerContentType, contentTypeJSON).
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
