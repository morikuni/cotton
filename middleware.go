package yacm

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func (m Middleware) WrapService(w http.ResponseWriter, r *http.Request, s Service) error {
	var err error
	h := m(http.HandlerFunc(func(w2 http.ResponseWriter, r2 *http.Request) {
		err = s.TryServeHTTP(w2, r2)
	}))
	h.ServeHTTP(w, r)
	return err
}
