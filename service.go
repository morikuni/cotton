package yacm

import (
	"net/http"
)

// Service is similar to http.Handler but it may returns an error.
type Service interface {
	// TryServeHTTP handles the request and respond to the client but it may fails with an error.
	TryServeHTTP(w http.ResponseWriter, r *http.Request) error
}

// ServiceFunc is the adapter to use a function as Service.
type ServiceFunc func(w http.ResponseWriter, r *http.Request) error

// TryServeHTTP implements Service.
func (f ServiceFunc) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
