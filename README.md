# poolit

Small, low overhead object pool to Go!

## Example

### Generic Version

```go
import "github.com/rzvxa/poolit"

func ExampleGenericPool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeGenericPool(
		10,
		func() *MyType { return new(MyType) }, // new
		func(mt *MyType) { *mt = MyType{} },   // cleanup
	)

	a := pool.Get()
	b := pool.Get()

	// use a and b

	pool.Release(a)
	pool.Release(b)
}
```

### Interface Version

###### Usually you want this for a mix of safety and convenient

```go
import "github.com/rzvxa/poolit"

func ExampleObjectPool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeObjectPool(
		10,
		func() any { return new(MyType) }, // new
		func(it any) {                     // cleanup
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
```

### Unsafe Version

In this version you are responsible to ensure the safety of `UnsafePool.Get` and `UnsafePool.Release` calls, But it comes with minimal indirection.

```go
import (
	"unsafe"

	"github.com/rzvxa/poolit"
)

func ExampleUnsafePool() {
	type MyType struct {
		value string
	}
	pool := poolit.MakeUnsafePool(
		10,
		func() unsafe.Pointer { return unsafe.Pointer(new(MyType)) }, // new
		func(p unsafe.Pointer) {                                      // cleanup
			*(*MyType)(p) = MyType{}
		},
	)

	a := (*MyType)(pool.Get())
	b := (*MyType)(pool.Get())

	// use a and b

	pool.Release(unsafe.Pointer(a))
	pool.Release(unsafe.Pointer(b))
}
```
