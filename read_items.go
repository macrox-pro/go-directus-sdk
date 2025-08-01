package directus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type ReadItemsPayload[T any] struct {
	Data []T `json:"data"`
}

type DeepQuery struct {
	Filter any      `json:"_filter,omitempty"`
	Sort   []string `json:"_sort,omitempty"`
	Search string   `json:"_search,omitempty"`
	Offset int      `json:"_offset,omitempty"`
	Limit  int      `json:"_limit,omitempty"`
	Page   int      `json:"_page,omitempty"`
}

type ReadItemsQuery struct {
	Alias  map[string]string    `url:"alias,omitempty"`
	Filter helpers.URLParamJSON `url:"filter,omitempty"`
	Sort   []string             `url:"sort,omitempty"`
	Deep   helpers.URLParamJSON `url:"deep,omitempty"`
	Search string               `url:"search,omitempty"`
	Offset int                  `url:"offset,omitempty"`
	Limit  int                  `url:"limit,omitempty"`
	Page   int                  `url:"page,omitempty"`
}

type ReadItemsRequest[T any] struct {
	ReadItemsQuery

	Collection string
	IsSystem   bool

	Token string

	ctx context.Context
}

func (r *ReadItemsRequest[T]) SetDeep(v map[string]DeepQuery) *ReadItemsRequest[T] {
	r.Deep = helpers.URLParamJSON{Data: v}
	return r
}

func (r *ReadItemsRequest[T]) SetAlias(v map[string]string) *ReadItemsRequest[T] {
	r.Alias = v
	return r
}

func (r *ReadItemsRequest[T]) SetFilter(rule FilterRule) *ReadItemsRequest[T] {
	r.Filter = helpers.URLParamJSON{Data: rule}
	return r
}

func (r *ReadItemsRequest[T]) SetSearch(search string) *ReadItemsRequest[T] {
	r.Search = search
	return r
}

func (r *ReadItemsRequest[T]) SetOffset(offset int) *ReadItemsRequest[T] {
	r.Offset = offset
	return r
}

func (r *ReadItemsRequest[T]) SetLimit(limit int) *ReadItemsRequest[T] {
	r.Limit = limit
	return r
}

func (r *ReadItemsRequest[T]) SetPage(page int) *ReadItemsRequest[T] {
	r.Page = page
	return r
}

func (r *ReadItemsRequest[T]) SetSort(sort []string) *ReadItemsRequest[T] {
	r.Sort = sort
	return r
}

func (r *ReadItemsRequest[T]) SetToken(token string) *ReadItemsRequest[T] {
	r.Token = token
	return r
}

func (r *ReadItemsRequest[T]) SetContext(ctx context.Context) *ReadItemsRequest[T] {
	r.ctx = ctx
	return r
}

func (r *ReadItemsRequest[T]) SetIsSystem(v bool) *ReadItemsRequest[T] {
	r.IsSystem = v
	return r
}

func (r *ReadItemsRequest[T]) SendBy(client *Client) ([]T, error) {
	if r.Collection == "" {
		return nil, fmt.Errorf("empty collection name")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	var payload ReadItemsPayload[T]

	req.QueryParam, _ = query.Values(r.ReadItemsQuery)
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
		return nil, err
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
		return nil, err
	}

	return payload.Data, nil
}

func NewReadItems[T any](collection string) *ReadItemsRequest[T] {
	return &ReadItemsRequest[T]{
		Collection: collection,
	}
}
