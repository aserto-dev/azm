package query

import set "github.com/deckarep/golang-set/v2"

func NewSet[T comparable](vals ...T) set.Set[T] {
	return set.NewThreadUnsafeSet(vals...)
}

func SetFromSlice[S any, T comparable](l []S, transform func(S) T) set.Set[T] {
	s := set.NewThreadUnsafeSetWithSize[T](len(l))
	for _, v := range l {
		s.Add(transform(v))
	}

	return s
}
