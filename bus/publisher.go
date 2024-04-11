package bus

import (
	"context"
)

type Publisher[T any] interface {
	Publish(ctx context.Context, messages ...Message[T]) error
}

type ClientPublisher[T any] interface {
	Publisher[T]
	Close(ctx context.Context) error
}
