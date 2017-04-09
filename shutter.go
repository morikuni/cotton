package yacm

import (
	"io"
	"net/http"
)

var (
	DefaultShutter = ShutterFunc(defaultShutter)
)

type Shutter interface {
	ShutError(w http.ResponseWriter, r *http.Request, err error)
}

type ShutterFunc func(w http.ResponseWriter, r *http.Request, err error)

func (f ShutterFunc) ShutError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

func defaultShutter(w http.ResponseWriter, _ *http.Request, _ error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, http.StatusText(http.StatusInternalServerError))
}
