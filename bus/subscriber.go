package bus

import (
	"context"
	"fmt"
)

type HandleStatus int32

const (
	StatusOk HandleStatus = iota
	StatusWarning
	StatusError
)

type HandleResult struct {
	Err  error
	Code HandleStatus
}

type Subscriber[T any] interface {
	Handle(ctx context.Context, message Message[T]) *HandleResult
}

var ErrClientClosed = fmt.Errorf("client closed")

type ClientConsumer interface {
	Consume(ctx context.Context) error
	Close(ctx context.Context) error
}
