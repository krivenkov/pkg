package es

import (
	"context"
	"encoding/json"

	"github.com/krivenkov/pkg/option"
	"github.com/krivenkov/pkg/order"
	"github.com/krivenkov/pkg/paginator"
	"github.com/olivere/elastic/v7"
)

type Client interface {
	Save(ctx context.Context, req *SaveRequest) error
	DeleteByID(ctx context.Context, index string, id string) error
	GetSearch(ctx context.Context, req *GetSearchRequest) (*GetSearchResponse, error)

	OpenPIT(ctx context.Context, indexes []string, keepAlive string) (*PIT, error)
	ClosePIT(ctx context.Context, pit *PIT) error
}

type SearchQuery interface {
	Source() (interface{}, error)
}

type RefreshType string

const (
	// RefreshTypeTrue refresh the relevant primary and replica shards immediately after the operation occurs
	RefreshTypeTrue RefreshType = "true"
	// RefreshTypeFalse the changes made by this request will be made visible at some point after the request returns
	RefreshTypeFalse RefreshType = "false"
	// RefreshTypeWaitFor elasticsearch automatically refreshes shards that have changed every index.refresh_interval which defaults to one second
	RefreshTypeWaitFor RefreshType = "wait_for"
)

type SaveRequest struct {
	Index   string
	ID      string
	Refresh RefreshType
	Obj     interface{}
}

type GetSearchRequest struct {
	Index         string
	Query         SearchQuery
	IncludeFields option.Option[[]string]
	ExcludeFields option.Option[[]string]
	MinScore      option.Option[float64]
	Highlight     option.Option[elastic.Highlight]
	Pit           option.Option[PIT]
	Orders        option.Option[[]*order.Order]
	Pagination    option.Option[paginator.Pagination]
}

type GetSearchResponse struct {
	Result []*SearchResult
	Total  int
}

type SearchResult struct {
	Source    json.RawMessage
	Highlight map[string][]string
}

type PIT struct {
	ID string `json:"id"`
}
