package middleware

import (
	"context"
	"fmt"
	router "github.com/mparaiso/simple-router-go"
	"net/http"
	"net/url"
	"time"
)

// StatusError is a status error
// it can be used to convert a http status to
// a go error interface
type StatusError int

func (se StatusError) Error() string {
	return http.StatusText(int(se))
}

// Container contains server values
type Container interface {
	ResponseWriter() http.ResponseWriter
	Request() *http.Request
}

// DefaultContainer is the default implementation of the Container
type DefaultContainer struct {
	RW  http.ResponseWriter // ResponseWriter
	Req *http.Request       // Request
}

// ResponseWriter returns a response writer
func (dc DefaultContainer) ResponseWriter() http.ResponseWriter { return dc.RW }

// Request returns a request
func (dc DefaultContainer) Request() *http.Request { return dc.Req }

// GetURLValues return URL variables
func (dc *DefaultContainer) GetURLValues() *url.Values {
	values := dc.Request().Context().Value(router.URLValues)
	if values == nil {
		dc.Req = dc.Req.WithContext(context.WithValue(dc.Request().Context(), router.URLValues, new(url.Values)))
		values = dc.Req.Context().Value(router.URLValues)
	}
	return values.(*url.Values)
}

// Error writes an error to the client and logs an error to stdout
func (dc *DefaultContainer) Error(statusCode int, err error) {
	http.Error(dc.ResponseWriter(), http.StatusText(statusCode), statusCode)
	fmt.Printf("[ERROR] %s [%s] \"%s %s %d\" : %s \n", dc.Request().RemoteAddr, time.Now().Format(""), dc.Request().Method, dc.Request().URL.RequestURI(), statusCode, err)
}

// Redirect replies with a redirection
func (dc *DefaultContainer) Redirect(url string, statusCode int) {
	http.Redirect(dc.ResponseWriter(), dc.Request(), url, statusCode)
}

// Handler is a controller that takes a context
type Handler func(Container)

// Wrap wraps Route.Handler with each middleware and returns a new Route
func (h Handler) Wrap(middlewares ...func(Handler) Handler) Handler {
	handler := h
	for i := len(middlewares) - 1; i == 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (h Handler) Handle(c Container) {
	h(c)
}

// ToHandlerFunc converts Handler to http.Handler
func (h Handler) ToHandlerFunc(containerFactory func(rw http.ResponseWriter, r *http.Request) Container) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		c := containerFactory(rw, r)
		h(c)
	}
}

// ToMiddleware wraps a classic net/http middleware (func(http.HandlerFunc) http.HandlerFunc)
// into a Middleware compatible with this package
func ToMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc) Middleware {
	return func(h Handler) Handler {
		return func(c Container) {
			middleware(h.ToHandlerFunc(func(rw http.ResponseWriter, r *http.Request) Container {
				return c
			})).ServeHTTP(c.ResponseWriter(), c.Request())
		}
	}
}

type Middleware func(Handler) Handler

func (m Middleware) Then(middleware Middleware) Middleware {
	return func(h Handler) Handler {
		return m(middleware(h))
	}
}

func (m Middleware) Finish(h Handler) Handler {
	return m(h)
}
