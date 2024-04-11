package builder

import (
	"context"
	"fmt"

	"github.com/krivenkov/pkg/bus"
	"github.com/krivenkov/pkg/bus/franz"
	"github.com/krivenkov/pkg/busapi/topics"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewPublisher[T any](cfg Config, logger *zap.Logger, topic topics.Topic) (bus.ClientPublisher[T], error) {
	switch cfg.Transport {
	case TransportFranz:
		return franz.NewPublisher[T](cfg.Franz, logger, topic)
	}

	return nil, fmt.Errorf("transport '%s' unsuported", cfg.Transport)
}

func NewFXPublisher[T any](topic topics.Topic) func(cfg Config, logger *zap.Logger, lc fx.Lifecycle) (bus.Publisher[T], error) {
	return func(cfg Config, logger *zap.Logger, lc fx.Lifecycle) (bus.Publisher[T], error) {
		cp, err := NewPublisher[T](cfg, logger, topic)
		if err != nil {
			return nil, fmt.Errorf("cannot create publisher for topic=[%s]: %w", topic, err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return cp.Close(ctx)
			},
		})

		return cp, nil
	}
}

func NewConsumer[T any](cfg Config, logger *zap.Logger, topic topics.Topic, groupID string, sub bus.Subscriber[T]) (bus.ClientConsumer, error) {
	return newConsumer(cfg, false, logger, topic, groupID, sub)
}

func newConsumer[T any](cfg Config, isOld bool, logger *zap.Logger, topic topics.Topic, groupID string, sub bus.Subscriber[T]) (bus.ClientConsumer, error) {
	switch cfg.Transport {
	case TransportFranz:
		return franz.NewConsumer[T](cfg.Franz, isOld, logger, topic, groupID, sub)
	}

	return nil, fmt.Errorf("transport '%s' unsuported", cfg.Transport)
}
