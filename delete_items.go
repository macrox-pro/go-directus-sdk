package directus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type DeleteItemsQuery struct {
	Filter any `url:"filter,omitempty"`
}

type DeleteItemsRequest[ID comparable] struct {
	DeleteItemsQuery

	Collection string
	IsSystem   bool

	IDs []ID

	Token string

	ctx context.Context
}

func (r *DeleteItemsRequest[ID]) SetIDs(ids ...ID) *DeleteItemsRequest[ID] {
	r.IDs = ids
	return r
}

func (r *DeleteItemsRequest[ID]) SetFilter(filter any) *DeleteItemsRequest[ID] {
	r.Filter = filter
	return r
}

func (r *DeleteItemsRequest[ID]) SetToken(token string) *DeleteItemsRequest[ID] {
	r.Token = token
	return r
}

func (r *DeleteItemsRequest[ID]) SetContext(ctx context.Context) *DeleteItemsRequest[ID] {
	r.ctx = ctx
	return r
}

func (r *DeleteItemsRequest[ID]) SetIsSystem(v bool) *DeleteItemsRequest[ID] {
	r.IsSystem = v
	return r
}

func (r *DeleteItemsRequest[ID]) SendBy(client *Client) error {
	if r.Collection == "" {
		return fmt.Errorf("empty collection name")
	}

	if len(r.IDs) == 0 && r.Filter == nil {
		return fmt.Errorf("empty delete conditions (id or filter)")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	if r.Filter != nil {
		req.QueryParam, _ = query.Values(r.DeleteItemsQuery)
	}

	if len(r.IDs) > 0 {
		req.Body = &r.IDs
	}

	resp, err := req.Delete(
		helpers.JoinPartsURL(
			helpers.PartURL{}, // for prefix - /
			helpers.PartURL{Value: "items", Skip: r.IsSystem},
			helpers.PartURL{Value: r.Collection},
		),
	)
	if err != nil {
		return err
	}

	body := resp.RawBody()
	if body == nil {
		return nil
	}

	defer body.Close()

	if !resp.IsError() {
		return nil
	}

	var payload DeleteItemPayload
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return err
	}

	return payload.Errors
}

func NewDeleteItems[ID comparable](collection string, ids ...ID) *DeleteItemsRequest[ID] {
	return &DeleteItemsRequest[ID]{
		Collection: collection,
		IDs:        ids,
	}
}
