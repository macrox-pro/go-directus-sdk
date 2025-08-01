package directus

import (
	"context"
	"encoding/json"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type AuthMode string

const (
	JsonAuthMode    AuthMode = "json"
	CookieAuthMode  AuthMode = "cookie"
	SessionAuthMode AuthMode = "session"
)

type AuthResult struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
}

type AuthResponsePayload struct {
	Data AuthResult `json:"data"`
}

type AuthLoginOptions struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Mode     AuthMode `json:"mode,omitempty"`
	OTP      string   `json:"otp,omitempty"`

	Provider string `json:"-"`
}

func (c *Client) AuthLogin(ctx context.Context, options AuthLoginOptions) (AuthResult, error) {
	var payload AuthResponsePayload

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Context-Type", "application/json").
		SetBody(options).
		Post(
			helpers.JoinPartsURL(
				helpers.PartURL{Value: "/auth/login"},
				helpers.PartURL{Value: options.Provider, Skip: options.Provider == ""},
			),
		)

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
