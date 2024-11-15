package poolit_test

import (
	"testing"

	"github.com/rzvxa/poolit"
)

func TestObjectPoolCapacity(t *testing.T) {
	newCount := 0
	expectNewCount := func(expect int) {
		if newCount != expect {
			t.Fatalf("Created more objects than initially requested, Expected %d, Created %d", expect, newCount)
		}
	}
	new := func() any {
		newCount++
		return new(Type)
	}
	pool := poolit.MakeObjectPool(2, new, nil)
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
func TestObjectPoolLIFO(t *testing.T) {
	new := func() any {
		return new(Type)
	}
	pool := poolit.MakeObjectPool(2, new, nil)

	objA1 := pool.Get()
	objB1 := pool.Get()

	pool.Release(objB1)
	pool.Release(objA1)

	objA2 := pool.Get()
	objB2 := pool.Get()

	if objA1 != objA2 || objB1 != objB2 {
		t.Fatalf("Expected `ObjectPool` to have a `LIFO` behavior")
	}
}

func TestObjectPoolCleanup(t *testing.T) {
	new := func() any {
		return &Type{c: "fresh"}
	}
	cleanCount := 0
	cleanup := func(it any) {
		cleanCount++
		x := it.(*Type)
		*x = Type{
			a: true,
			b: 42,
			c: "clean",
		}
	}
	pool := poolit.MakeObjectPool(1, new, cleanup)

	obj := pool.Get().(*Type)

	if obj.a || obj.b != 0 || obj.c != "fresh" {
		t.Fatal("Invalid data")
	}

	obj.c = "dirty"

	pool.Release(obj)

	obj = pool.Get().(*Type)

	// expect it to be clean
	if obj.c != "clean" {
		t.Fatalf("Invalid data. Expected `obj.c` to be `clean`, Found %s", obj.c)
	}

	if cleanCount != 1 {
		t.Fatalf("Expected cleanup to be called once, But it's called %d times", cleanCount)
	}
}

// new items should be fresh(aka don't go through a redundant cleanup)
func TestObjectPoolFreshness(t *testing.T) {
	new := func() any {
		return &Type{c: "fresh"}
	}
	cleanCount := 0
	cleanup := func(it any) {
		cleanCount++
		x := it.(*Type)
		*x = Type{c: "clean"}
	}
	pool := poolit.MakeObjectPool(1, new, cleanup)

	pool.Get()
	obj := pool.Get().(*Type)

	// expect it to be clean
	if obj.c != "fresh" {
		t.Fatalf("Invalid data. Expected `obj.c` to be `fresh`, Found %s", obj.c)
	}

	if cleanCount != 0 {
		t.Fatalf("Expected cleanup to never be called, But it's called %d times", cleanCount)
	}
}
