package directus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"

	"github.com/macrox-pro/go-directus-sdk/helpers"
)

type AggregateQuery struct {
	Aggregate AggregateRule        `url:"aggregate,omitempty"`
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

func (r *AggregateRequest[T]) SetAggregate(rule AggregateRule) *AggregateRequest[T] {
	r.Aggregate = rule
	return r
}

func (r *AggregateRequest[T]) SetGroupBy(groupBy []string) *AggregateRequest[T] {
	r.GroupBy = groupBy
	return r
}

func (r *AggregateRequest[T]) SetFilter(rule FilterRule) *AggregateRequest[T] {
	r.Filter = helpers.URLParamJSON{Data: rule}
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
	var payload ReadItemsPayload[T]

	if r.Collection == "" {
		return payload.Data, fmt.Errorf("empty collection name")
	}

	if len(r.GroupBy) == 0 && r.Aggregate == nil {
		return payload.Data, fmt.Errorf("invalid request aggregate options")
	}

	req := client.createRequestWithContext(r.ctx).
		SetDoNotParseResponse(true).
		SetHeader(headerAccept, contentTypeJSON)

	if r.Token != "" {
		req.SetAuthToken(r.Token)
	}

	req.QueryParam, _ = query.Values(r.AggregateQuery)

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
		var failed ErrorResponse
		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
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

func NewAggregate[T any](collection string) *AggregateRequest[T] {
	return &AggregateRequest[T]{
		Collection: collection,
	}
}
