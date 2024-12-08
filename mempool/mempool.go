package mempool

import "sync"

const defaultSliceCapacity = 128

type Pool[T any] struct {
	sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}

func NewPool[T any](newF func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() interface{} {
				return newF()
			},
		},
	}
}

func NewSlicePool[T any]() *Pool[*[]T] {
	return NewPool(func() *[]T {
		s := make([]T, 0, defaultSliceCapacity)
		return &s
	})
}

type Resetable[T any] interface {
	Reset()
	*T
}

type CollectionPool[M any, T Resetable[M]] struct {
	slicePool *Pool[*[]T]
	msgPool   *Pool[T]
}

func NewCollectionPool[M any, T Resetable[M]]() *CollectionPool[M, T] {
	return &CollectionPool[M, T]{
		slicePool: NewSlicePool[T](),
		msgPool: NewPool(func() T {
			return new(M)
		}),
	}
}

func (p CollectionPool[M, T]) GetSlice() *[]T {
	// return p.slicePool.Get()
	return p.slicePool.New().(*[]T)
}

func (p *CollectionPool[M, T]) PutSlice(s *[]T) {
	for _, item := range *s {
		item.Reset()
		p.msgPool.Put(item)
	}

	*s = (*s)[:0]
	p.slicePool.Put(s)
}

func (p *CollectionPool[M, T]) Get() T {
	// return p.msgPool.Get()
	return p.msgPool.New().(T)
}

func (p *CollectionPool[M, T]) Put(m T) {
	p.msgPool.Put(m)
}
