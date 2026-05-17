package directus

import (
	"context"
	"encoding/json"
)

func (c *Client) ServerPing(ctx context.Context) (string, error) {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader(headerAccept, contentTypeJSON).
		Get("/server/ping")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		var failed ErrorsPayload
		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
			return "", err
		}
		return "", failed.Errors
	}

	var payload struct {
		Ping string `json:"ping"`
	}

	err = json.Unmarshal(resp.Body(), &payload)
	return payload.Ping, err
}
