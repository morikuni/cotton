package yacm

import (
	"net/http"
)

// Service is a extended http.HandlerFunc.
// Error may be returned (nil on success).
type Service func(w http.ResponseWriter, r *http.Request) Error

// Recover registers a ErrorHandler as a Error handler.
// Registered function is called only when the Service returned a Error.
func (s Service) Recover(h ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s(w, r)
		if err != nil {
			h(w, r, err)
		}
	}
}

// IgnoreError ignores a Error of the Service.
func (s Service) IgnoreError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s(w, r)
	}
}
