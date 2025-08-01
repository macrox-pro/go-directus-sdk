package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type ReadItemQuery struct {
	Deep helpers.URLParamJSON `url:"deep,omitempty"`
}

type ReadItemPayload[T any] struct {
	Data T `json:"data"`
}

type ReadItemRequest[T any] struct {
	ReadItemQuery

	Collection, ID string
	IsSystem       bool

	Token string

	ctx context.Context
}

func (r *ReadItemRequest[T]) SetDeep(v map[string]DeepQuery) *ReadItemRequest[T] {
	r.Deep = helpers.URLParamJSON{Data: v}
	return r
}

func (r *ReadItemRequest[T]) SetToken(token string) *ReadItemRequest[T] {
	r.Token = token
	return r
}

func (r *ReadItemRequest[T]) SetContext(ctx context.Context) *ReadItemRequest[T] {
	r.ctx = ctx
	return r
}

func (r *ReadItemRequest[T]) SetIsSystem(v bool) *ReadItemRequest[T] {
	r.IsSystem = v
	return r
}

func (r *ReadItemRequest[T]) SendBy(client *Client) (T, error) {
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

	req.QueryParam, _ = query.Values(r.ReadItemQuery)
	if req.QueryParam == nil {
		req.QueryParam = url.Values{}
	}
	req.QueryParam["fields"] = helpers.ExtractFieldsJSON(payload.Data)

	resp, err := req.Get(
		helpers.JoinPartsURL(
			helpers.PartURL{}, // for prefix - /
			helpers.PartURL{Value: "items", Skip: r.IsSystem},
			helpers.PartURL{Value: r.Collection},
			helpers.PartURL{Value: r.ID},
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

func NewReadItem[T any](collection, id string) *ReadItemRequest[T] {
	return &ReadItemRequest[T]{
		Collection: collection,
		ID:         id,
	}
}
