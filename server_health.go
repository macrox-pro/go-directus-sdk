package directus

import (
	"context"
	"encoding/json"
)

func (c *Client) ServerHealth(ctx context.Context) (string, error) {
	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get("/server/health")

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
		Status string `json:"status"`
	}

	err = json.Unmarshal(resp.Body(), &payload)
	return payload.Status, err
}
