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

func Compose(middlewares ...Middleware) Middleware {
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
		next := Compose(middlewares[1:]...)
		return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
			Apply(m, Apply(next, h)).ServeHTTP(w, r)
		})
	}
}

func Apply(m Middleware, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.WrapHandler(w, r, h)
	})
}

func MiddlewareToFilter(m Middleware) Filter {
	return FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		var err error
		m.WrapHandler(w, r, http.HandlerFunc(func(w2 http.ResponseWriter, r2 *http.Request) {
			err = s.TryServeHTTP(w2, r2)
		}))
		return err
	})
}
