package poolit

import "log"

type ObjectPool struct {
	items   []any
	new     anyNewFn
	cleanup anyCleanupFn
	inuse   int
}

type anyNewFn = func() any
type anyCleanupFn = func(any)

func MakeObjectPool(initialSize int, new anyNewFn, cleanup anyCleanupFn) ObjectPool {
	if new == nil {
		log.Fatalln("The `new` function should never be `nil`")
	}
	if cleanup == nil {
		cleanup = func(a any) {}
	}
	p := ObjectPool{
		items:   make([]any, initialSize),
		new:     new,
		cleanup: cleanup,
	}
	for ix := 0; ix < initialSize; ix++ {
		p.items[ix] = new()
	}
	return p
}

func (p *ObjectPool) Get() any {
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

func (p *ObjectPool) Release(it any) {
	p.cleanup(it)
	p.inuse--
	ix := p.ix()
	p.items[ix] = it
}

func (p *ObjectPool) ix() int {
	return len(p.items) - p.inuse - 1
}
