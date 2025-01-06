package ds

import (
	"iter"
	"maps"

	"github.com/aserto-dev/azm/mempool"
	"github.com/samber/lo"
)

type nothing struct{}

type Set[T comparable] struct {
	items map[T]nothing
}

func NewSet[T comparable](capacity int) Set[T] {
	return Set[T]{make(map[T]nothing, capacity)}
}

func (s Set[T]) Cardinality() int {
	return len(s.items)
}

func (s Set[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

func (s Set[T]) Contains(v T) (ok bool) {
	_, ok = s.items[v]
	return ok
}

func (s Set[T]) Elements() iter.Seq[T] {
	if s.items == nil {
		return EmptyIter[T]()
	}
	return maps.Keys(s.items)
}

func (s Set[T]) ToSlice() []T {
	return lo.Keys(s.items)
}

func (s *Set[T]) Add(vals iter.Seq[T]) {
	for val := range vals {
		s.add(val)
	}
}

func (s *Set[T]) Remove(vals iter.Seq[T]) {
	for v := range vals {
		delete(s.items, v)
	}
}

func (s *Set[T]) Union(other Set[T]) {
	s.Add(other.Elements())
}

func (s *Set[T]) Intersect(other Set[T]) {
	s.Remove(FilterIter(s.Elements(), func(v T) bool {
		return !other.Contains(v)
	}))
}

func (s *Set[T]) Difference(other Set[T]) {
	s.Remove(FilterIter(s.Elements(), other.Contains))
}

func (s Set[T]) add(v T) {
	(s.items)[v] = nothing{}
}

func (s Set[T]) removeAll() {
	for k := range s.items {
		delete(s.items, k)
	}
}

type SetPool[T comparable] struct {
	pool *mempool.Pool[map[T]nothing]
}

func NewSetPool[T comparable]() *SetPool[T] {
	return &SetPool[T]{
		mempool.NewPool(func() map[T]nothing {
			return make(map[T]nothing)
		})}
}

func (p *SetPool[T]) GetSet() Set[T] {
	return Set[T]{p.pool.Get()}
}

func (p *SetPool[T]) PutSet(s Set[T]) {
	if s.items == nil {
		return
	}

	s.removeAll()
	p.pool.Put(s.items)
	s.items = nil
}
