package yacm

import (
	"io"
	"net/http"
)

var (
	// DefaultShutter handles all errors as 500 internal server error.
	DefaultShutter = ShutterFunc(defaultShutter)
)

// Shutter handles an error from Filter or Service.
type Shutter interface {
	// ShutError converts the error to response.
	ShutError(w http.ResponseWriter, r *http.Request, err error)
}

// ShutterFunc is the adapter to user a function as Shutter.
type ShutterFunc func(w http.ResponseWriter, r *http.Request, err error)

// ShutError implements Shutter.
func (f ShutterFunc) ShutError(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

func defaultShutter(w http.ResponseWriter, _ *http.Request, _ error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, http.StatusText(http.StatusInternalServerError))
}
