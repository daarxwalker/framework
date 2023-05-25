package framework

type Signal[T any] struct {
	value T
}

func (s *Signal[T]) Update(updater func(value T) T) {
	s.value = updater(s.value)
}

func (s *Signal[T]) Set(value T) {
	s.value = value
}

func (s *Signal[T]) Get() T {
	return s.value
}

func (s *Signal[T]) Watch() T {
	return s.value
}
