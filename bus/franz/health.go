package franz

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type health struct {
	cfg Config
}

func NewHealth(cfg Config) (*health, error) {
	return &health{cfg: cfg}, nil
}

func (h *health) Check(ctx context.Context) error {
	opts := []kgo.Opt{
		kgo.SeedBrokers(h.cfg.Addresses...),
	}

	cli, err := kgo.NewClient(opts...)
	if err != nil {
		return err
	}

	defer cli.Close()

	return cli.Ping(ctx)
}
