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
	return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		for i := len(middlewares) - 1; i >= 0; i-- {
			m := middlewares[i]
			h = Apply(m, h)
		}
		h.ServeHTTP(w, r)
	})
}

func Apply(m Middleware, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.WrapHandler(w, r, h)
	})
}

// // ToFilter changes Middleware into a Filter.
// func (m Middleware) ToFilter() Filter {
// 	return func(w http.ResponseWriter, r *http.Request, s Service) error {
// 		var err error
// 		m(w, r, func(w2 http.ResponseWriter, r2 *http.Request) {
// 			err = s(w2, r2)
// 		})
// 		return err
// 	}
// }
