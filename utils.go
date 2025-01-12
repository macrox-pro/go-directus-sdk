package directus

import (
	"context"
	"encoding/json"
	"strconv"
)

func (c *Client) RandomString(ctx context.Context, length int) (string, error) {
	resp, err := c.createRequestWithContext(ctx).
		SetHeader(headerAccept, contentTypeJSON).
		SetQueryParam("length", strconv.Itoa(length)).
		Get("/utils/random/string")
	if err != nil {
		return "", err
	}

	var payload ReadItemPayload[string]

	err = json.Unmarshal(resp.Body(), &payload)
	if err != nil {
		return "", err
	}

	if resp.IsError() && len(payload.Errors) > 0 {
		return payload.Data, payload.Errors
	}

	return payload.Data, nil
}
