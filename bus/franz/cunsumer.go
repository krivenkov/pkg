package franz

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/krivenkov/pkg/bus"
	"github.com/krivenkov/pkg/busapi/topics"
	"github.com/krivenkov/pkg/mlog"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type consumer[T any] struct {
	cli    *kgo.Client
	sub    bus.Subscriber[T]
	logger *zap.Logger
	isOld  bool
}

func NewConsumer[T any](cfg Config, isOld bool, logger *zap.Logger, topic topics.Topic, groupID string, sub bus.Subscriber[T]) (*consumer[T], error) {
	logger = logger.With(
		zap.String("busProvider", "FRANZ"),
		zap.String("busSubTopic", string(topic)),
		zap.String("busGroupID", groupID),
	)

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Addresses...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(string(topic)),
		kgo.HeartbeatInterval(cfg.HeartbeatInterval),
		kgo.FetchMaxWait(cfg.MaxWait),
		kgo.FetchMaxBytes(cfg.FetchMaxBytes),
		kgo.BrokerMaxReadBytes(cfg.BrokerMaxReadBytes),
		kgo.SessionTimeout(cfg.SessionTimeout),
		kgo.DisableAutoCommit(),
	}
	cli, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	logger.Info("consumer created")

	return &consumer[T]{cli: cli, sub: sub, logger: logger, isOld: isOld}, nil
}

func (c *consumer[T]) Consume(ctx context.Context) error {
	for {
		fetches := c.cli.PollFetches(ctx)

		if fetches.IsClientClosed() {
			return bus.ErrClientClosed
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, fErr := range errs {
				c.logger.Error("fetches error",
					zap.Error(fErr.Err),
					zap.Int32("partition", fErr.Partition),
					zap.String("topic", fErr.Topic),
				)
			}
		}

		for iter := fetches.RecordIter(); !iter.Done(); {
			c.handleRecord(ctx, iter.Next())
		}
	}
}

func (c *consumer[T]) Close(_ context.Context) error {
	c.cli.Close()
	c.logger.Info("consumer closed")
	return nil
}

func (c *consumer[T]) handleRecord(ctx context.Context, kr *kgo.Record) {
	logger := c.logger
	ctx = mlog.CtxWithLogger(ctx, logger)

	logger.Debug("start record handle")

	defer func() {
		if rec := recover(); rec != nil {
			logger.Error("recover handle record",
				zap.Any("rec", rec),
				zap.Stack("stack"),
			)
		}
	}()

	var payload T

	val := bus.MessageValue[T]{
		CreatedAt: time.Time{},
		Payload:   &payload,
	}

	prv := reflect.ValueOf(&payload)
	isProto, _ := implementsInterface(prv.Type(), protoMessageType)

	switch {
	case isProto:
		if err := proto.Unmarshal(kr.Value, prv.Interface().(proto.Message)); err != nil {
			logger.Error("unmarshall proto record error", zap.Error(err))
			return
		}
	case c.isOld:
		if err := json.Unmarshal(kr.Value, &payload); err != nil {
			logger.Error("unmarshall json record error", zap.Error(err))
			return
		}
	default:
		if err := json.Unmarshal(kr.Value, &val); err != nil {
			logger.Error("unmarshall record error", zap.Error(err))
			return
		}
	}

	res := c.sub.Handle(ctx, bus.Message[T]{
		Key:   string(kr.Key),
		Value: val,
	})

	switch res.Code {
	case bus.StatusError:
		logger.Error("handle record error", zap.Error(res.Err))
		return
	case bus.StatusWarning:
		logger.Warn("handle record warn", zap.Error(res.Err))
	}

	if err := c.cli.CommitRecords(ctx, kr); err != nil {
		logger.Error("commit record error", zap.Error(err))
	}
}

var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

func implementsInterface(typ, decType reflect.Type) (success bool, indir int8) {
	if typ == nil {
		return
	}
	rt := typ

	for {
		if rt.Implements(decType) {
			return true, indir
		}
		if p := rt; p.Kind() == reflect.Pointer {
			indir++
			if indir > 100 {
				return false, 0
			}
			rt = p.Elem()
			continue
		}
		break
	}
	if typ.Kind() != reflect.Pointer {
		if reflect.PointerTo(typ).Implements(decType) {
			return true, -1
		}
	}

	return false, 0
}
