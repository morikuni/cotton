package yacm

import (
	"net/http"
)

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error) error
}

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error) error

func (f ErrorHandlerFunc) HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	return f(w, r, err)
}

type chainedErrorHandler []ErrorHandler

func (c chainedErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	for _, h := range c {
		err = h.HandleError(w, r, err)
		if err == nil {
			return nil
		}
	}
}

func ComposeErrorHandler(handlers ...ErrorHandler) ErrorHandler {
	l := len(handlers)
	switch l {
	case 0:
		panic(ErrEmptyArgs)
	default:
		return chainedErrorHandler(handlers)
	}
}
