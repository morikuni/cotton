package yacm

import (
	"net/http"
)

type Service interface {
	TryServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type ServiceFunc func(w http.ResponseWriter, r *http.Request) error

func (f ServiceFunc) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
