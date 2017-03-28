package yacm

import (
	"net/http"
)

type ErrorShutter interface {
	ShutError(w http.ResponseWriter, r *http.Request, err error)
}

type ErrorShutterFunc func(w http.ResponseWriter, r *http.Request, err error)

func (f ErrorShutterFunc) ShutError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}
