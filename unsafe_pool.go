package poolit

import (
	"log"
	"unsafe"
)

// SAFETY: all items used with the `UnsafePool` should have the same type,
// The user is responsible for keeping pointer conversions type-safe
type UnsafePool struct {
	items   []unsafe.Pointer
	new     unsafeNewFn
	cleanup unsafeCleanupFn
	inuse   int
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
		items:   make([]unsafe.Pointer, initialSize),
		new:     new,
		cleanup: cleanup,
	}
	for ix := 0; ix < initialSize; ix++ {
		p.items[ix] = new()
	}
	return p
}

// SAFETY: Caller is responsible to cast the `unsafe.Pointer` to the correct type
func (p *UnsafePool) Get() unsafe.Pointer {
	ix := p.ix()
	p.inuse++

	// if we are out of objects resize the pool
	if ix == -1 {
		// each time we run out add 2 more items to the pool
		// one for this call and one extra item reserved for later calls
		p.items = append(p.items, nil, nil)
		p.items[ix+1] = p.new() // extra item
		return p.new()
	}
	return p.items[ix]
}

// SAFETY: Caller is responsible to only release appropiate types which are expected by the callers of `UnsafePool.Get` method
func (p *UnsafePool) Release(ptr unsafe.Pointer) {
	p.cleanup(ptr)
	p.inuse--
	ix := p.ix()
	p.items[ix] = ptr
}

func (p *UnsafePool) ix() int {
	return len(p.items) - p.inuse - 1
}
