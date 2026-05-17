package directus

import "context"

var accessTokenContextKey = struct{}{}

func WithAccessTokenContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, token)
}

func AccessTokenFromContext(ctx context.Context) string {
	if i := ctx.Value(accessTokenContextKey); i != nil {
		if token, ok := i.(string); ok {
			return token
		}
	}

	return ""
}
