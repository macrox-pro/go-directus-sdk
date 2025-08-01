package directus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type DeleteItemRequest struct {
	Collection, ID string
	IsSystem       bool

	Token string

	ctx context.Context
}

func (r *DeleteItemRequest) SetToken(token string) *DeleteItemRequest {
	r.Token = token
	return r
}

func (r *DeleteItemRequest) SetContext(ctx context.Context) *DeleteItemRequest {
	r.ctx = ctx
	return r
}

func (r *DeleteItemRequest) SetIsSystem(v bool) *DeleteItemRequest {
	r.IsSystem = v
	return r
}

func (r *DeleteItemRequest) SendBy(client *Client) error {
	if r.Collection == "" {
		return fmt.Errorf("empty collection name")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	resp, err := req.Delete(
		helpers.JoinPartsURL(
			helpers.PartURL{}, // for prefix - /
			helpers.PartURL{Value: "items", Skip: r.IsSystem},
			helpers.PartURL{Value: r.Collection},
			helpers.PartURL{Value: r.ID},
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

func NewDeleteItem(collection, id string) *DeleteItemRequest {
	return &DeleteItemRequest{
		Collection: collection,
		ID:         id,
	}
}
