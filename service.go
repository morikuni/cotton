package cotton

import (
	"net/http"
)

// Service is a extended http.HandlerFunc.
// Error may be returned (nil on success).
type Service func(w http.ResponseWriter, r *http.Request) Error

// Recover registers a RecoverFunc as a Error handler.
// Registered function is called only when the Service returned a Error.
func (s Service) Recover(f RecoverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s(w, r)
		if err != nil {
			f(w, r, err)
		}
	}
}
