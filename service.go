package yacm

import (
	"net/http"
)

// Service is a extended http.HandlerFunc.
type Service func(w http.ResponseWriter, r *http.Request) error

// Recover registers a ErrorHandler as a error handler.
// Registered function is called only when the Service returned a error.
func (s Service) Recover(h ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s(w, r)
		if err != nil {
			h(w, r, err)
		}
	}
}

// IgnoreError ignores a error of the Service.
func (s Service) IgnoreError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s(w, r)
	}
}
