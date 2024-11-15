package poolit_test

import (
	"log"
	"unsafe"

	"github.com/rzvxa/poolit"
)

func ExampleGenericPool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeGenericPool(
		10,
		func() *MyType { return new(MyType) },
		func(mt *MyType) { *mt = MyType{} },
	)

	a := pool.Get()
	b := pool.Get()

	// use a and b

	pool.Release(a)
	pool.Release(b)
}

func ExampleObjectPool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeObjectPool(
		10,
		func() any { return new(MyType) },
		func(it any) {
			val, ok := it.(*MyType)
			if !ok {
				log.Fatalln("Invalid object")
			}
			*val = MyType{}
		},
	)

	a := pool.Get().(*MyType)
	b := pool.Get().(*MyType)

	// use a and b

	pool.Release(a)
	pool.Release(b)
}

func ExampleUnsafePool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeUnsafePool(
		10,
		func() unsafe.Pointer { return unsafe.Pointer(new(MyType)) },
		func(p unsafe.Pointer) {
			*(*MyType)(p) = MyType{}
		},
	)

	a := (*MyType)(pool.Get())
	b := (*MyType)(pool.Get())

	// use a and b

	pool.Release(unsafe.Pointer(a))
	pool.Release(unsafe.Pointer(b))
}
