package yacm

import (
	"io"
	"net/http"
)

var (
	DefaultErrorShutter = ErrorShutterFunc(defaultErrorShutter)
)

type ErrorShutter interface {
	ShutError(w http.ResponseWriter, r *http.Request, err error)
}

type ErrorShutterFunc func(w http.ResponseWriter, r *http.Request, err error)

func (f ErrorShutterFunc) ShutError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

func defaultErrorShutter(w http.ResponseWriter, _ *http.Request, _ error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, http.StatusText(http.StatusInternalServerError))
}
