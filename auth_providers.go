package directus

import (
	"context"
	"encoding/json"
)

type AuthProvider struct {
	Name        string            `json:"name"`
	Driver      string            `json:"driver"`
	Icon        string            `json:"icon,omitempty"`
	Label       string            `json:"label,omitempty"`
	RedirectURL string            `json:"redirect_url,omitempty"`
	AuthURL     string            `json:"auth_url,omitempty"`
	ClientID    string            `json:"client_id,omitempty"`
	Scope       []string          `json:"scope,omitempty"`
	Additional  map[string]string `json:"additional,omitempty"`
}

func (c *Client) AuthProviders(ctx context.Context) ([]AuthProvider, error) {
	var payload ReadItemsPayload[AuthProvider]

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get("/auth/providers")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		var failed ErrorsPayload
		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
			return nil, err
		}
		return nil, failed.Errors
	}

	err = json.Unmarshal(resp.Body(), &payload)
	return payload.Data, err
}
