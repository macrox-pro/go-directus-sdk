package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type CreateItemRequest[T any] struct {
	ReadItemQuery

	Collection string
	IsSystem   bool
	Data       any

	Token string

	ctx context.Context
}

func (r *CreateItemRequest[T]) SetData(data any) *CreateItemRequest[T] {
	r.Data = data
	return r
}

func (r *CreateItemRequest[T]) SetDeep(v map[string]DeepQuery) *CreateItemRequest[T] {
	r.Deep = helpers.URLParamJSON{Data: v}
	return r
}

func (r *CreateItemRequest[T]) SetToken(token string) *CreateItemRequest[T] {
	r.Token = token
	return r
}

func (r *CreateItemRequest[T]) SetContext(ctx context.Context) *CreateItemRequest[T] {
	r.ctx = ctx
	return r
}

func (r *CreateItemRequest[T]) SetIsSystem(v bool) *CreateItemRequest[T] {
	r.IsSystem = v
	return r
}

func (r *CreateItemRequest[T]) SendBy(client *Client) (T, error) {
	var payload ReadItemPayload[T]

	if r.Collection == "" {
		return payload.Data, fmt.Errorf("empty collection name")
	}

	if r.Data == nil {
		return payload.Data, fmt.Errorf("empty data")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerContentType, contentTypeJSON).
		SetHeader(headerAccept, contentTypeJSON).
		SetBody(r.Data)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	req.QueryParam, _ = query.Values(r.ReadItemQuery)
	if req.QueryParam == nil {
		req.QueryParam = url.Values{}
	}
	req.QueryParam["fields"] = helpers.ExtractFieldsJSON(payload.Data)

	resp, err := req.Post(
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

func NewCreateItem[T any](collection string, data any) *CreateItemRequest[T] {
	return &CreateItemRequest[T]{
		Collection: collection,
		Data:       data,
	}
}
