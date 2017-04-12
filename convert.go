package yacm

import (
	"net/http"
)

func ServiceToHandler(s Service, shutter Shutter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.TryServeHTTP(w, r)
		if err != nil {
			shutter.ShutError(w, r, err)
		}
	})
}

func HandlerToService(h http.Handler) Service {
	return ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	})
}
