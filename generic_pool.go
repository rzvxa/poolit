package poolit

import "unsafe"

type GenericPool[T any] struct {
	pool UnsafePool
}

func makeGenericPool[T any](initialSize int, new func() *T, cleanup func(*T)) GenericPool[T] {
	return GenericPool[T]{
		pool: MakeUnsafePool(
			initialSize,
			func() unsafe.Pointer { return unsafe.Pointer(new()) },
			func(p unsafe.Pointer) { cleanup((*T)(p)) },
		),
	}
}

func (p *GenericPool[T]) Get() *T {
	return (*T)(p.pool.Get())
}

func (p *GenericPool[T]) Release(it *T) {
	p.pool.Release(unsafe.Pointer(it))
}
