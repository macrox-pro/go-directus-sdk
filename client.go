package directus

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"

	contentTypeJSON = "application/json"

	defaultAuthScheme = "Bearer"
)

type Client struct {
	resty *resty.Client

	extractTokenFromContext bool
}

func (c *Client) createRequestWithContext(ctx context.Context) *resty.Request {
	if ctx == nil {
		return c.resty.R()
	}

	restyReq := c.resty.R().SetContext(ctx)

	if !c.extractTokenFromContext {
		return restyReq
	}

	if token := AccessTokenFromContext(ctx); token != "" {
		restyReq.SetAuthToken(token)
	}

	return restyReq
}

type ClientOption func(c *Client)

func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, err
	}

	restyClient := resty.New()
	restyClient.BaseURL = strings.TrimRight(baseURL, "/")
	restyClient.AuthScheme = defaultAuthScheme

	client := &Client{resty: restyClient}
	for _, fn := range options {
		fn(client)
	}

	return client, nil
}
