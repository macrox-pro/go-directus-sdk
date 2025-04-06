package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type UpdateItemsRequetBody struct {
	Data any      `json:"data"`
	Keys []string `json:"keys,omitempty"`
}

type UpdateItemsRequest[T any] struct {
	ReadItemsQuery

	Collection string
	Keys       []string

	IsSystem bool

	Changes any

	Token string

	ctx context.Context
}

func (r *UpdateItemsRequest[T]) SetChanges(changes any) *UpdateItemsRequest[T] {
	r.Changes = changes
	return r
}

func (r *UpdateItemsRequest[T]) SetDeep(v map[string]DeepQuery) *UpdateItemsRequest[T] {
	r.Deep = helpers.URLParamJSON{Data: v}
	return r
}

func (r *UpdateItemsRequest[T]) SetFilter(rule FilterRule) *UpdateItemsRequest[T] {
	r.Filter = helpers.URLParamJSON{Data: rule}
	return r
}

func (r *UpdateItemsRequest[T]) SetToken(token string) *UpdateItemsRequest[T] {
	r.Token = token
	return r
}

func (r *UpdateItemsRequest[T]) SetContext(ctx context.Context) *UpdateItemsRequest[T] {
	r.ctx = ctx
	return r
}

func (r *UpdateItemsRequest[T]) SetIsSystem(v bool) *UpdateItemsRequest[T] {
	r.IsSystem = v
	return r
}

func (r *UpdateItemsRequest[T]) SendBy(client *Client) (T, error) {
	var payload ReadItemPayload[T]

	if r.Collection == "" {
		return payload.Data, fmt.Errorf("empty collection name")
	}

	if r.Changes == nil {
		return payload.Data, fmt.Errorf("empty changes")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerContentType, contentTypeJSON).
		SetHeader(headerAccept, contentTypeJSON).
		SetBody(&UpdateItemsRequetBody{
			Data: r.Changes,
			Keys: r.Keys,
		})

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	req.QueryParam, _ = query.Values(r.ReadItemsQuery)
	if req.QueryParam == nil {
		req.QueryParam = url.Values{}
	}
	req.QueryParam["fields"] = helpers.ExtractFieldsJSON(payload.Data)

	resp, err := req.Patch(
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

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload.Data, err
	}

	if resp.IsError() && len(payload.Errors) > 0 {
		return payload.Data, payload.Errors
	}

	return payload.Data, nil
}

func NewUpdateItems[T any](collection string, keys []string, changes any) *UpdateItemsRequest[T] {
	return &UpdateItemsRequest[T]{
		Collection: collection,
		Changes:    changes,
		Keys:       keys,
	}
}
