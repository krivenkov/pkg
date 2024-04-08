package txer

import "context"

//go:generate mockgen -source=txer.go -package=txer_mock -destination=mock/txer.go

type TXer interface {
	WithTX(context.Context, func(context.Context) error) error
}
