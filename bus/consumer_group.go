package bus

import (
	"context"

	"github.com/oklog/run"
)

type ConsumerGroup struct {
	consumers []ClientConsumer
}

func NewConsumerGroup(consumers ...ClientConsumer) *ConsumerGroup {
	return &ConsumerGroup{consumers: consumers}
}

func (cg *ConsumerGroup) Add(consumer ClientConsumer) {
	cg.consumers = append(cg.consumers, consumer)
}

func (cg *ConsumerGroup) Consume(ctx context.Context) error {
	g := &run.Group{}

	for _, c := range cg.consumers {
		func(c ClientConsumer) {
			g.Add(func() error {
				return c.Consume(ctx)
			}, func(_ error) {
				c.Close(ctx)
			})
		}(c)
	}

	return g.Run()
}

func (cg *ConsumerGroup) Close(ctx context.Context) error {
	for _, c := range cg.consumers {
		if err := c.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}
