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

	if resp.IsError() {
		var failed ErrorResponse

		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
			return "", Error{
				Status:  resp.StatusCode(),
				Details: err,
			}
		}

		return "", Error{
			Status:  resp.StatusCode(),
			Details: failed.Errors,
		}
	}

	var payload ReadItemPayload[string]

	err = json.Unmarshal(resp.Body(), &payload)
	if err != nil {
		return "", err
	}

	return payload.Data, nil
}
