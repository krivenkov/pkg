// Package franz contains implementation publisher and subscriber for franz-go client
package franz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/krivenkov/pkg/bus"
	"github.com/krivenkov/pkg/busapi/topics"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type pub[T any] struct {
	cli    *kgo.Client
	topic  topics.Topic
	logger *zap.Logger
}

func NewPublisher[T any](cfg Config, logger *zap.Logger, topic topics.Topic) (*pub[T], error) {
	logger = logger.With(
		zap.String("busProvider", "FRANZ"),
		zap.String("busPubTopic", string(topic)),
	)

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Addresses...),
		kgo.ProducerLinger(cfg.BatchTimeout),
		kgo.HeartbeatInterval(cfg.HeartbeatInterval),
		kgo.SessionTimeout(cfg.SessionTimeout),
		kgo.BrokerMaxReadBytes(cfg.BrokerMaxReadBytes),
	}
	cli, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("publisher created")

	return &pub[T]{cli: cli, topic: topic, logger: logger}, nil
}

func (p *pub[T]) Publish(ctx context.Context, messages ...bus.Message[T]) error {
	records := make([]*kgo.Record, 0, len(messages))
	for _, message := range messages {
		data, err := json.Marshal(message.Value)
		if err != nil {
			return fmt.Errorf("value marshal to json: %w", err)
		}
		record := &kgo.Record{Topic: string(p.topic), Value: data}

		if len(message.Key) > 0 {
			record.Key = []byte(message.Key)
		}

		records = append(records, record)
	}

	if err := p.cli.ProduceSync(ctx, records...).FirstErr(); err != nil {
		p.logger.Error("error ProduceSync",
			zap.String("TOPIC", p.topic.String()),
			zap.Error(err))
		return err
	}

	p.logger.Debug("published successfully")

	return nil
}

func (p *pub[T]) Close(_ context.Context) error {
	p.cli.Close()
	p.logger.Info("publisher closed")
	return nil
}
