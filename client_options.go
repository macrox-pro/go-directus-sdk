package directus

func WithStaticToken(token string) ClientOption {
	return func(client *Client) {
		client.resty.SetAuthToken(token)
	}
}

func WithExtractTokenFromContext(enabled bool) ClientOption {
	return func(client *Client) {
		client.extractTokenFromContext = enabled
	}
}
