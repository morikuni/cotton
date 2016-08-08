package yacm

import (
	"net/http"
)

// Middleware is a middleware for http.HandlerFunc.
type Middleware func(w http.ResponseWriter, r *http.Request, s Service) error

// And composes two Middlewares.
func (m Middleware) And(next Middleware) Middleware {
	return func(w http.ResponseWriter, r *http.Request, s Service) error {
		return m(w, r, next.Then(s))
	}
}

// For wraps a given http.HandlerFunc and upgrades to Serice.
func (m Middleware) For(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		return m(w, r, handlerToService(handler))
	}
}

// Then wraps a given Service.
func (m Middleware) Then(s Service) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		return m(w, r, s)
	}
}

func handlerToService(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		handler(w, r)
		return nil
	}
}
