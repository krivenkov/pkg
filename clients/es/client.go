package es

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/krivenkov/pkg/mlog"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type client struct {
	*elastic.Client
}

func NewClient(cfg Config, logger *zap.Logger) (Client, error) {
	cli, err := elastic.NewClient(
		elastic.SetURL(cfg.Addresses...),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(time.Second*5),
		elastic.SetHttpClient(http.DefaultClient),
		elastic.SetTraceLog(newLoggerES(logger)),
	)
	if err != nil {
		return nil, err
	}
	return &client{
		Client: cli,
	}, nil
}

func (c *client) Save(ctx context.Context, req *SaveRequest) error {
	logger := mlog.FromContext(ctx).With(
		zap.String("index", req.Index),
		zap.String("id", req.ID),
	)

	if _, err := c.Index().
		Index(req.Index).
		BodyJson(req.Obj).
		Id(req.ID).
		Refresh(string(req.Refresh)).
		Do(ctx); err != nil {
		logger.Error("es fail create doc",
			zap.Error(err),
		)
		return err
	}

	logger.Debug("es create")

	return nil
}

func (c *client) DeleteByID(ctx context.Context, index string, id string) error {
	logger := mlog.FromContext(ctx).With(
		zap.String("index", index),
		zap.String("id", id),
	)

	if _, err := c.Delete().Index(index).Id(id).Refresh("true").Do(ctx); err != nil {
		if errObj, ok := err.(*elastic.Error); ok {
			if errObj.Status == http.StatusNotFound {
				logger.Warn("object not found", zap.Error(err))
				return ErrNotFound
			}

			return errObj
		}

		logger.Error("es fail delete doc", zap.Error(err))
	}

	logger.Debug("es delete")

	return nil
}

func (c *client) GetSearch(ctx context.Context, req *GetSearchRequest) (*GetSearchResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request in elastic search")
	}

	q := c.Client.Search()

	q.Index(req.Index)
	q.Query(req.Query)

	fsc := elastic.NewFetchSourceContext(true)

	if req.IncludeFields.IsSet() {
		fsc.Include(req.IncludeFields.Value()...)
	}

	if req.ExcludeFields.IsSet() {
		fsc.Exclude(req.ExcludeFields.Value()...)
	}

	if req.MinScore.IsSet() {
		q.MinScore(req.MinScore.Value())
	}

	q.FetchSourceContext(fsc)

	if req.Highlight.IsSet() {
		hl := req.Highlight.Value()
		q.Highlight(&hl)
	}

	if req.Pit.IsSet() {
		q.PointInTime(elastic.NewPointInTime(req.Pit.Value().ID))
	}

	if req.Pagination.IsSet() {
		q.From(req.Pagination.Value().Offset)
		q.Size(req.Pagination.Value().Limit)
	}

	for _, val := range req.Orders.Value() {
		if strings.ToLower(val.Direction) == "desc" {
			q.Sort(val.Column, false)
			continue
		}

		q.Sort(val.Column, true)
	}

	searchResult, err := q.TrackTotalHits(true).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*SearchResult, 0)
	for _, hit := range searchResult.Hits.Hits {
		res = append(res, &SearchResult{
			Source:    hit.Source,
			Highlight: hit.Highlight,
		})
	}

	return &GetSearchResponse{
		Result: res,
		Total:  int(searchResult.Hits.TotalHits.Value),
	}, nil
}

func (c *client) OpenPIT(ctx context.Context, indexes []string, keepAlive string) (*PIT, error) {
	logger := mlog.FromContext(ctx).With(
		zap.String("indexes", strings.Join(indexes, ",")),
	)

	openResp, err := c.Client.
		OpenPointInTime(indexes...).
		KeepAlive(keepAlive).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if openResp == nil {
		logger.Error("expected non-nil Point In Time")
		return nil, fmt.Errorf("expected non-nil Point In Time")
	}
	if openResp.Id == "" {
		logger.Error("expected non-blank Point In Time ID")
		return nil, fmt.Errorf("expected non-blank Point In Time ID")
	}

	return &PIT{
		ID: openResp.Id,
	}, nil
}

func (c *client) ClosePIT(ctx context.Context, pit *PIT) error {
	logger := mlog.FromContext(ctx).With(
		zap.String("pit_id", pit.ID),
	)

	closeResp, err := c.Client.ClosePointInTime(pit.ID).Pretty(true).Do(context.Background())
	if err != nil {
		logger.Error("can not close pit", zap.Error(err))
		return fmt.Errorf("can not close pit")
	}

	if closeResp == nil {
		logger.Error("can not close pit", zap.Error(err))
		return fmt.Errorf("can not close pit")
	}

	if !closeResp.Succeeded {
		logger.Error("can not close pit")
		return fmt.Errorf("can not close pit")
	}
	if closeResp.NumFreed != 1 {
		logger.Error("can not close pit")
		return fmt.Errorf("can not close pit")
	}

	return nil
}
