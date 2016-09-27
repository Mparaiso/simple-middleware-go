package middleware_test

import (
	"fmt"
	m "github.com/mparaiso/simple-middleware-go"
	"net/http"
	"net/http/httptest"
)

func ExampleMiddleware_Then() {

	// Let's chain middlewares thanks to the Then method

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
		Finish(func(m.Container) { fmt.Println("Handle the request") }).
		Handle(nil)

	// Output:
	// 123Handle the request
}

func ExampleToMiddleware() {
	// Let's convert a classic http middleware into a middleware supported by this package

	classicCORSMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			// this middleware handles corss origin requests from browsers
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(rw, r)
		}
	}

	convertedMiddleware := m.ToMiddleware(classicCORSMiddleware)

	// Let's test our converted middleware
	request, _ := http.NewRequest("GET", "https://acme.com", nil)
	response := httptest.NewRecorder()

	convertedMiddleware.
		Finish(func(c m.Container) { c.ResponseWriter().Write([]byte("done")) }).
		Handle(&m.DefaultContainer{response, request})

	fmt.Println(response.Header().Get("Access-Control-Allow-Origin"))
	fmt.Println(response.Body.String())

	// Output:
	// *
	// done

}
