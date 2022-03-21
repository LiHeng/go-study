package generic

type Ordered interface {
	int | string | int64 | int32 | int8 | int16 | float32 | float64 | uint8 | uint32 | uint | uint64 | uint16
}

func min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

type Tree[T any] struct {
	left, right *Tree[T]
	data        T
}

func (t *Tree[T]) Lookup(x T) *Tree[T] {
	return nil
}

type Set[T Ordered] struct {
	data map[T]any
}

func NewSet[T Ordered]() *Set[T] {
	return &Set[T]{
		data: make(map[T]any),
	}
}

func (s *Set[T]) Add(item T) {
	s.data[item] = struct{}{}
}

func (s *Set[T]) Values() []T {
	var items []T
	for key := range s.data {
		items = append(items, key)
	}
	return items
}

var stringTree Tree[string]
