package ptr

func Pointer[T any](value T) *T {
	newValue := value
	return &newValue
}

func CopyPointer[T any](value *T) *T {
	if value == nil {
		return nil
	}

	newValue := *value
	return &newValue
}
