package yacm

import (
	"net/http"
)

// Middleware is a middleware for http.HandlerFunc.
type Middleware func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc)

// Compose composes two Middlewares.
func (m Middleware) Compose(next Middleware) Middleware {
	return func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		m(w, r, next.Apply(h))
	}
}

// Apply wraps a given http.HandlerFunc.
func (m Middleware) Apply(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m(w, r, h)
	}
}

// ToFilter changes Middleware into a Filter.
func (m Middleware) ToFilter() Filter {
	return func(w http.ResponseWriter, r *http.Request, s Service) error {
		var err error
		m(w, r, func(w2 http.ResponseWriter, r2 *http.Request) {
			err = s(w2, r2)
		})
		return err
	}
}
