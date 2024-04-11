package bus

import "context"

type Health interface {
	Check(ctx context.Context) error
}
