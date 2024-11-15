package poolit

import (
	"testing"
	"unsafe"
)

type Type struct {
	a bool
	b int64
	c string
}

func TestCapacity(t *testing.T) {
	newCount := 0
	expectNewCount := func(expect int) {
		if newCount != expect {
			t.Fatalf("Created more objects than initially requested, Expected %d, Created %d", expect, newCount)
		}
	}
	new := func() unsafe.Pointer {
		newCount++
		return unsafe.Pointer(new(Type))
	}
	pool := MakeUnsafePool(2, new, nil)
	expectNewCount(2)

	// taking existing items shouldn't affect the pool size
	pool.Get()
	pool.Get()
	expectNewCount(2)

	// when we run out of items, we add 2 more at once
	pool.Get()
	expectNewCount(4)

	// we create no more instances as long as we are recycling items
	last := pool.Get()
	pool.Release(last)
	pool.Get()
	expectNewCount(4)
}

// Last In First Out
func TestLIFO(t *testing.T) {
	new := func() unsafe.Pointer {
		return unsafe.Pointer(new(Type))
	}
	pool := MakeUnsafePool(2, new, nil)

	objA1 := pool.Get()
	objB1 := pool.Get()

	pool.Release(objB1)
	pool.Release(objA1)

	objA2 := pool.Get()
	objB2 := pool.Get()

	if objA1 != objA2 || objB1 != objB2 {
		t.Fatalf("Expected `UnsafePool` to have a `LIFO` behavior")
	}
}

// Make sure to run tests with `-gcflags=all=-d=checkptr` to test pointers correctly
func TestSafety(t *testing.T) {
	new := func() unsafe.Pointer {
		return unsafe.Pointer(&Type{a: true, b: 42, c: "fresh"})
	}
	pool := MakeUnsafePool(2, new, nil)

	obj1 := (*Type)(pool.Get())
	obj2 := (*Type)(pool.Get())

	if !obj1.a || obj1.b != 42 || obj1.c != "fresh" {
		t.Fatalf("Invalid data")
	}

	obj2.c = "dirty"

	pool.Release(unsafe.Pointer(obj2))
	pool.Release(unsafe.Pointer(obj1))

	obj1 = (*Type)(pool.Get())
	obj2 = (*Type)(pool.Get())

	if !obj1.a || obj1.b != 42 || obj1.c != "fresh" || obj2.c != "dirty" {
		t.Fatalf("Invalid data")
	}
}

func TestCleanup(t *testing.T) {
	new := func() unsafe.Pointer {
		return unsafe.Pointer(&Type{c: "fresh"})
	}
	cleanCount := 0
	cleanup := func(ptr unsafe.Pointer) {
		cleanCount++
		*(*Type)(ptr) = Type{
			a: true,
			b: 42,
			c: "clean",
		}
	}
	pool := MakeUnsafePool(1, new, cleanup)

	obj := (*Type)(pool.Get())

	if obj.a || obj.b != 0 || obj.c != "fresh" {
		t.Fatal("Invalid data")
	}

	obj.c = "dirty"

	pool.Release(unsafe.Pointer(obj))

	obj = (*Type)(pool.Get())

	// expect it to be clean
	if obj.c != "clean" {
		t.Fatalf("Invalid data. Expected `obj.c` to be `clean`, Found %s", obj.c)
	}

	if cleanCount != 1 {
		t.Fatalf("Expected cleanup to be called once, But it's called %d times", cleanCount)
	}
}

// new items should be fresh(aka don't go through a redundant cleanup)
func TestFreshness(t *testing.T) {
	new := func() unsafe.Pointer {
		return unsafe.Pointer(&Type{c: "fresh"})
	}
	cleanCount := 0
	cleanup := func(ptr unsafe.Pointer) {
		cleanCount++
		*(*Type)(ptr) = Type{c: "clean"}
	}
	pool := MakeUnsafePool(1, new, cleanup)

	pool.Get()
	obj := (*Type)(pool.Get())

	// expect it to be clean
	if obj.c != "fresh" {
		t.Fatalf("Invalid data. Expected `obj.c` to be `fresh`, Found %s", obj.c)
	}

	if cleanCount != 0 {
		t.Fatalf("Expected cleanup to never be called, But it's called %d times", cleanCount)
	}
}
