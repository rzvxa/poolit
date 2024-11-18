package poolit

import (
	"unsafe"
)

// This is similar to `UnsafeThinPool` but with no events to have a truly indirection free object pool
// The `new` function is supplied as a parameter to the `UnsafeThinPool.Get` method so the go compiler have a chance of inlining it
type UnsafeThinPool struct {
	items []unsafe.Pointer
	inuse int
}

func MakeUnsafeThinPool(initialSize int, new unsafeNewFn) UnsafeThinPool {
	p := UnsafeThinPool{
		items: make([]unsafe.Pointer, initialSize),
	}
	for ix := 0; ix < initialSize; ix++ {
		p.items[ix] = new()
	}
	return p
}

// SAFETY: Caller is responsible to cast the `unsafe.Pointer` to the correct type
func (p *UnsafeThinPool) Get(new unsafeNewFn) unsafe.Pointer {
	ix := p.ix()
	p.inuse++

	// if we are out of objects resize the pool
	if ix == -1 {
		// each time we run out add 2 more items to the pool
		// one for this call and one extra item reserved for later calls
		p.items = append(p.items, nil, nil)
		p.items[ix+1] = new() // extra item
		return new()
	}
	return p.items[ix]
}

// SAFETY: Caller is responsible to only release appropiate types which are expected by the callers of `UnsafeThinPool.Get` method
func (p *UnsafeThinPool) Release(ptr unsafe.Pointer) {
	p.inuse--
	ix := p.ix()
	p.items[ix] = ptr
}

// number of inuse objects
func (p *UnsafeThinPool) InUse() int {
	return p.inuse
}

func (p *UnsafeThinPool) ix() int {
	return len(p.items) - p.inuse - 1
}
