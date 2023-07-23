package config

type fValue[T any] struct {
	value T
}

func (f *fValue[T]) String() T {
	return f.value
}

func (f *fValue[T]) Set(s T) error {
	f.value = s
	return nil
}

func (f *fValue[T]) Get() *T {
	return &f.value
}
