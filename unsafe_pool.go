package poolit

import (
	"log"
	"unsafe"
)

// SAFETY: all items used with the `UnsafePool` should have the same type,
// The user is responsible for keeping pointer conversions type-safe
type UnsafePool struct {
	pool    UnsafeThinPool
	new     unsafeNewFn
	cleanup unsafeCleanupFn
}

type unsafeNewFn = func() unsafe.Pointer
type unsafeCleanupFn = func(unsafe.Pointer)

func MakeUnsafePool(initialSize int, new unsafeNewFn, cleanup unsafeCleanupFn) UnsafePool {
	if new == nil {
		log.Fatalln("The `new` function should never be `nil`")
	}
	if cleanup == nil {
		cleanup = func(p unsafe.Pointer) {}
	}
	p := UnsafePool{
		pool:    MakeUnsafeThinPool(initialSize, new),
		new:     new,
		cleanup: cleanup,
	}
	return p
}

func NewUnsafePool(initialSize int, new unsafeNewFn, cleanup unsafeCleanupFn) *UnsafePool {
	self := MakeUnsafePool(initialSize, new, cleanup)
	return &self
}

// SAFETY: Caller is responsible to cast the `unsafe.Pointer` to the correct type
func (p *UnsafePool) Get() unsafe.Pointer {
	return p.pool.Get(p.new)
}

// SAFETY: Caller is responsible to only release appropiate types which are expected by the callers of `UnsafePool.Get` method
func (p *UnsafePool) Release(ptr unsafe.Pointer) {
	p.cleanup(ptr)
	p.pool.Release(ptr)
}
