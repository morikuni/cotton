package cotton

import (
	"net/http"
)

type Error interface{}

type Filter func(w http.ResponseWriter, r *http.Request, s Service) Error

func (f Filter) And(next Filter) Filter {
	return func(w http.ResponseWriter, r *http.Request, s Service) Error {
		return f(w, r, next.Then(s))
	}
}

func (f Filter) For(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		return f(w, r, handlerToService(handler))
	}
}

func (f Filter) Then(s Service) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		return f(w, r, s)
	}
}

func handlerToService(handler http.HandlerFunc) Service {
	return func(w http.ResponseWriter, r *http.Request) Error {
		handler(w, r)
		return nil
	}
}
