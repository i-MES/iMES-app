package imes

import (
	"fmt"
	"testing"
)

func TestAddCounter(t *testing.T) {
	mid := new(Middleware)
	i := mid.LoadCounter()
	if (i + 1) != mid.AddCounter() {
		t.Errorf("AddCounter from %d want %d", i, i+1)
	}
}

func BenchmarkAddCounter(b *testing.B) {
	mid := new(Middleware)
	for i := 0; i < b.N; i++ {
		mid.AddCounter()
	}
}

func BenchmarkAddCounterParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		mid := new(Middleware)
		for i := 0; i < b.N; i++ {
			mid.AddCounter()
		}
	})
}

func ExampleMiddleware_AddCounter() {
	mid := new(Middleware)
	mid.AddCounter()
	// output:
	// i+1
}

func ExampleMiddleware_LoadTestitems() {
	fmt.Println("foobar")
	mid := new(Middleware)
	mid.LoadTestitems("foobar")
	// Output:
	// {}
}
