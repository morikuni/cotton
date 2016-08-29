package yacm

import (
	"net/http"
)

// Filter is a middleware for http.HandlerFunc with error.
type Filter func(w http.ResponseWriter, r *http.Request, s Service) error

// Compose composes two Filters.
func (f Filter) Compose(next Filter) Filter {
	return func(w http.ResponseWriter, r *http.Request, s Service) error {
		return f(w, r, next.Then(s))
	}
}

// For wraps a given http.HandlerFunc and upgrades to Serice.
func (f Filter) For(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		return f(w, r, handlerToService(handler))
	}
}

// Then wraps a given Service.
func (f Filter) Then(s Service) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		return f(w, r, s)
	}
}

// Recover registers a ErrorHandler as a error handler.
// Registered function is called only when the Filter returned a error.
func (f Filter) Recover(eh ErrorHandler) Middleware {
	return func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		err := f(w, r, handlerToService(h))
		if err != nil {
			eh(w, r, err)
		}
	}
}

func handlerToService(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) error {
		handler(w, r)
		return nil
	}
}
