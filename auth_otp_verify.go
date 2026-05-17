package directus

import (
	"context"
	"encoding/json"
)

type OTPVerifyParams struct {
	OTP string `json:"otp"`
}

func (c *Client) AuthOTPVerify(ctx context.Context, options OTPVerifyParams) (AuthResult, error) {
	var payload AuthResponsePayload

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader(headerContentType, contentTypeJSON).
		SetBody(options).
		Post("/auth/otp/verify")

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
