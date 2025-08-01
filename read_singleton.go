package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type ReadSingletonQuery struct {
	Deep map[string]DeepQuery `url:"deep,omitempty"`
}

type ReadSingletonRequest[T any] struct {
	ReadSingletonQuery

	Collection string
	IsSystem   bool

	Token string

	ctx context.Context
}

func (r *ReadSingletonRequest[T]) SetDeep(deep map[string]DeepQuery) *ReadSingletonRequest[T] {
	r.Deep = deep
	return r
}

func (r *ReadSingletonRequest[T]) SetToken(token string) *ReadSingletonRequest[T] {
	r.Token = token
	return r
}

func (r *ReadSingletonRequest[T]) SetContext(ctx context.Context) *ReadSingletonRequest[T] {
	r.ctx = ctx
	return r
}

func (r *ReadSingletonRequest[T]) SetIsSystem(v bool) *ReadSingletonRequest[T] {
	r.IsSystem = v
	return r
}

func (r *ReadSingletonRequest[T]) SendBy(client *Client) (T, error) {
	var payload ReadItemPayload[T]

	if r.Collection == "" {
		return payload.Data, fmt.Errorf("empty collection name")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	req.QueryParam, _ = query.Values(r.ReadSingletonQuery)
	if req.QueryParam == nil {
		req.QueryParam = url.Values{}
	}
	req.QueryParam["fields"] = helpers.ExtractFieldsJSON(payload.Data)

	resp, err := req.Get(
		helpers.JoinPartsURL(
			helpers.PartURL{}, // for prefix - /
			helpers.PartURL{Value: "items", Skip: r.IsSystem},
			helpers.PartURL{Value: r.Collection},
		),
	)
	if err != nil {
		return payload.Data, err
	}

	body := resp.RawBody()
	if body != nil {
		defer body.Close()
	}

	if resp.IsError() {
		var failed ErrorsPayload
		if err := json.NewDecoder(body).Decode(&failed); err != nil {
			return payload.Data, err
		}

		return payload.Data, failed.Errors
	}

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload.Data, err
	}

	return payload.Data, nil
}

func NewReadSingleton[T any](collection string) *ReadSingletonRequest[T] {
	return &ReadSingletonRequest[T]{
		Collection: collection,
	}
}
