package middleware_test

import (
	"fmt"
	m "github.com/mparaiso/simple-middleware-go"
)

func ExampleMiddleware_Then() {
	middleware1 := func(next m.Handler) m.Handler {
		return func(c m.Container) { fmt.Print(1); next(c) }
	}
	middleware2 := func(next m.Handler) m.Handler {
		return func(c m.Container) { fmt.Print(2); next(c) }
	}
	middleware3 := func(next m.Handler) m.Handler {
		return func(c m.Container) { fmt.Print(3); next(c) }
	}
	m.Middleware(middleware1).
		Then(middleware2).
		Then(middleware3).
		Finish(func(m.Container) { fmt.Printf("Finish") }).
		Handle(nil)

	// Output:
	// 123Finish
}
