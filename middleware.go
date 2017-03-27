package yacm

import (
	"net/http"
)

type Middleware interface {
	WrapHandler(w http.ResponseWriter, r *http.Request, h http.Handler)
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, h http.Handler)

func (m MiddlewareFunc) WrapHandler(w http.ResponseWriter, r *http.Request, h http.Handler) {
	m(w, r, h)
}

func ComposeMiddleware(middlewares ...Middleware) Middleware {
	l := len(middlewares)
	switch l {
	case 0:
		return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
			h.ServeHTTP(w, r)
		})
	case 1:
		return middlewares[0]
	default:
		m := middlewares[0]
		next := ComposeMiddleware(middlewares[1:]...)
		return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
			ApplyMiddleware(m, ApplyMiddleware(next, h)).ServeHTTP(w, r)
		})
	}
}

func ApplyMiddleware(m Middleware, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.WrapHandler(w, r, h)
	})
}
