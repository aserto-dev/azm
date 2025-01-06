package ds

import "iter"

func EmptyIter[T any]() iter.Seq[T] {
	return func(_ func(T) bool) {}
}

func TransformIter[S any, T any](src iter.Seq[S], transform func(S) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range src {
			if !yield(transform(v)) {
				break
			}
		}
	}
}

func FilterIter[T any](src iter.Seq[T], filter func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range src {
			if filter(v) {
				if !yield(v) {
					break
				}
			}
		}
	}
}
