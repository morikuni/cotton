package yacm

import (
	"net/http"
)

// Middleware is a middleware for http.HandlerFunc.
type Middleware func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc)

// Compose composes two Middlewares.
func (m Middleware) Compose(next Middleware) Middleware {
	return func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		m(w, r, next.For(h))
	}
}

// For wraps a given http.HandlerFunc.
func (m Middleware) For(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m(w, r, h)
	}
}
