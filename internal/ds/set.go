package ds

import (
	"iter"
	"maps"

	"github.com/samber/lo"
)

type nothing struct{}

type Set[T comparable] map[T]nothing

func NewSet[T comparable](vals ...T) Set[T] {
	return SetFromSlice(vals, func(v T) T { return v })
}

func SetFromSlice[S any, T comparable](l []S, transform func(S) T) Set[T] {
	s := make(Set[T], len(l))
	for _, v := range l {
		s.Add(transform(v))
	}

	return s
}

func (s Set[T]) Cardinality() int {
	return len(s)
}

func (s Set[T]) IsEmpty() bool {
	return s.Cardinality() == 0
}

func (s Set[T]) Add(v T) bool {
	prevLen := len(s)
	s.add(v)
	return prevLen != len(s)
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	n := s.Cardinality()
	if other.Cardinality() > n {
		n = other.Cardinality()
	}
	unionedSet := make(Set[T], n)

	for elem := range s {
		unionedSet.add(elem)
	}
	for elem := range other {
		unionedSet.add(elem)
	}
	return unionedSet
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	intersection := NewSet[T]()

	// loop over smaller set
	var (
		smaller Set[T]
		larger  Set[T]
	)
	if s.Cardinality() < other.Cardinality() {
		smaller = s
		larger = other
	} else {
		smaller = other
		larger = s
	}

	for elem := range smaller {
		if larger.contains(elem) {
			intersection.add(elem)
		}
	}
	return intersection
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	diff := NewSet[T]()
	for elem := range s {
		if !other.contains(elem) {
			diff.add(elem)
		}
	}
	return diff
}

func (s Set[T]) Elements() iter.Seq[T] {
	return maps.Keys(s)
}

func (s Set[T]) ToSlice() []T {
	return lo.Keys(s)
}

func (s Set[T]) add(v T) {
	s[v] = nothing{}
}

func (s Set[T]) contains(v T) (ok bool) {
	_, ok = s[v]
	return ok
}
