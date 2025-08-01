package directus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type DeleteItemsQuery struct {
	Filter FilterRule `json:"filter"`
}

type DeleteItemsRequestPayload struct {
	Query DeleteItemsQuery `json:"query"`
}

type DeleteItemsRequest[ID comparable] struct {
	DeleteItemsRequestPayload

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

func (r *DeleteItemsRequest[ID]) SetFilter(rule FilterRule) *DeleteItemsRequest[ID] {
	r.Query.Filter = rule
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

	if len(r.IDs) == 0 && r.Query.Filter == nil {
		return fmt.Errorf("empty delete conditions (id or filter)")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON).
		SetHeader(headerContentType, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	if len(r.IDs) > 0 {
		req.Body = &r.IDs
	} else if r.Query.Filter != nil {
		req.Body = r.DeleteItemsRequestPayload
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

	if resp.IsSuccess() {
		return nil
	}

	var payload ErrorsPayload
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
