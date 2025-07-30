package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type UpdateItemRequest[T any] struct {
	ReadItemQuery

	Collection, ID string
	IsSystem       bool
	Changes        any

	Token string

	ctx context.Context
}

func (r *UpdateItemRequest[T]) SetChanges(changes any) *UpdateItemRequest[T] {
	r.Changes = changes
	return r
}

func (r *UpdateItemRequest[T]) SetDeep(v map[string]DeepQuery) *UpdateItemRequest[T] {
	r.Deep = helpers.URLParamJSON{Data: v}
	return r
}

func (r *UpdateItemRequest[T]) SetToken(token string) *UpdateItemRequest[T] {
	r.Token = token
	return r
}

func (r *UpdateItemRequest[T]) SetContext(ctx context.Context) *UpdateItemRequest[T] {
	r.ctx = ctx
	return r
}

func (r *UpdateItemRequest[T]) SetIsSystem(v bool) *UpdateItemRequest[T] {
	r.IsSystem = v
	return r
}

func (r *UpdateItemRequest[T]) SendBy(client *Client) (T, error) {
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
		SetBody(r.Changes)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	req.QueryParam, _ = query.Values(r.ReadItemQuery)
	if req.QueryParam == nil {
		req.QueryParam = url.Values{}
	}
	req.QueryParam["fields"] = helpers.ExtractFieldsJSON(payload.Data)

	resp, err := req.Patch(
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
		var failed ErrorResponse
		if err := json.NewDecoder(body).Decode(&failed); err != nil {
			return payload.Data, err
		}

		return payload.Data, Error{
			Status:  resp.StatusCode(),
			Details: failed.Errors,
		}
	}

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload.Data, err
	}

	return payload.Data, nil
}

func NewUpdateItem[T any](collection, id string, changes any) *UpdateItemRequest[T] {
	return &UpdateItemRequest[T]{
		Collection: collection,
		Changes:    changes,
		ID:         id,
	}
}
