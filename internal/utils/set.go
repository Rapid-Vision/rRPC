package utils

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}
