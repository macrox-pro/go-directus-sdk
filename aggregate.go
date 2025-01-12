package directus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type Count struct {
	Fields []string `url:"count,omitempty"`
}

type CountDistinct struct {
	Fields []string `url:"countDistinct,omitempty"`
}

type Sum struct {
	Fields []string `url:"sum,omitempty"`
}

type SumDistinct struct {
	Fields []string `url:"sumDistinct,omitempty"`
}

type Avg struct {
	Fields []string `url:"avg,omitempty"`
}

type AvgDistinct struct {
	Fields []string `url:"avgDistinct,omitempty"`
}

type Min struct {
	Field string `url:"min,omitempty"`
}

type Max struct {
	Field string `url:"max,omitempty"`
}

type AggregateQuery struct {
	Aggregate any                  `url:"aggregate,omitempty"`
	GroupBy   []string             `url:"groupBy,omitempty"`
	Filter    helpers.URLParamJSON `url:"filter,omitempty"`
	Sort      []string             `url:"sort,omitempty"`
	Search    string               `url:"search,omitempty"`
	Offset    int                  `url:"offset,omitempty"`
	Limit     int                  `url:"limit,omitempty"`
	Page      int                  `url:"page,omitempty"`
}

type AggregateRequest[T any] struct {
	AggregateQuery

	Collection string
	IsSystem   bool

	Token string

	ctx context.Context
}

func (r *AggregateRequest[T]) SetAggregate(aggregate any) *AggregateRequest[T] {
	r.Aggregate = aggregate
	return r
}

func (r *AggregateRequest[T]) SetGroupBy(groupBy []string) *AggregateRequest[T] {
	r.GroupBy = groupBy
	return r
}

func (r *AggregateRequest[T]) SetFilter(v any) *AggregateRequest[T] {
	r.Filter = helpers.URLParamJSON{Data: v}
	return r
}

func (r *AggregateRequest[T]) SetSearch(search string) *AggregateRequest[T] {
	r.Search = search
	return r
}

func (r *AggregateRequest[T]) SetOffset(offset int) *AggregateRequest[T] {
	r.Offset = offset
	return r
}

func (r *AggregateRequest[T]) SetLimit(limit int) *AggregateRequest[T] {
	r.Limit = limit
	return r
}

func (r *AggregateRequest[T]) SetPage(page int) *AggregateRequest[T] {
	r.Page = page
	return r
}

func (r *AggregateRequest[T]) SetSort(sort []string) *AggregateRequest[T] {
	r.Sort = sort
	return r
}

func (r *AggregateRequest[T]) SetToken(token string) *AggregateRequest[T] {
	r.Token = token
	return r
}

func (r *AggregateRequest[T]) SetContext(ctx context.Context) *AggregateRequest[T] {
	r.ctx = ctx
	return r
}

func (r *AggregateRequest[T]) SetIsSystem(v bool) *AggregateRequest[T] {
	r.IsSystem = v
	return r
}

func (r *AggregateRequest[T]) SendBy(client *Client) ([]T, error) {
	if r.Collection == "" {
		return nil, fmt.Errorf("empty collection name")
	}

	if len(r.GroupBy) == 0 && r.Aggregate == nil {
		return nil, fmt.Errorf("invalid request aggregate options")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	var payload ReadItemsPayload[T]

	req.QueryParam, _ = query.Values(r.AggregateQuery)

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

	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, err
	}

	if resp.IsError() && len(payload.Errors) > 0 {
		return nil, payload.Errors
	}

	return payload.Data, nil
}

func NewAggregate[T any](collection string) *AggregateRequest[T] {
	return &AggregateRequest[T]{
		Collection: collection,
	}
}
