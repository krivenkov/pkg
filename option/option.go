package option

import "encoding/json"

type Option[T any] struct {
	value *T
}

func New[T any](v T) Option[T] {
	return Option[T]{&v}
}

func Nil[T any]() Option[T] {
	return Option[T]{}
}

func (o *Option[T]) Set(v T) {
	o.value = &v
}

func (o Option[T]) IsSet() bool {
	return o.value != nil
}

func (o Option[T]) Value() T {
	if !o.IsSet() {
		var zero T
		return zero
	}
	return *o.value
}

func (o Option[T]) ValuePtr() *T {
	if o.IsSet() {
		newValue := *o.value
		return &newValue
	}
	return nil
}

func (o Option[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = &value
	return nil
}
