package cotton

import (
	"net/http"
)

type Error interface{}

type Middleware func(w http.ResponseWriter, r *http.Request, s Service) Error

func (m Middleware) And(next Middleware) Middleware {
	return func(w http.ResponseWriter, r *http.Request, s Service) Error {
		return m(w, r, next.Then(s))
	}
}

func (m Middleware) For(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		return m(w, r, handlerToService(handler))
	}
}

func (m Middleware) Then(s Service) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		return m(w, r, s)
	}
}

func handlerToService(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		handler(w, r)
		return nil
	}
}
