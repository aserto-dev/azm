package ds

type Stack[T any] struct {
	items []T
}

func NewStack[T any](items ...T) *Stack[T] {
	return &Stack[T]{items}
}

func (s Stack[T]) Len() int {
	return len(s.items)
}

func (s Stack[T]) Top() T {
	return s.items[s.Len()-1]
}

func (s Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	item := s.items[s.Len()-1]
	s.items = s.items[:s.Len()-1]
	return item
}
