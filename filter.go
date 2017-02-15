package yacm

import (
	"net/http"
)

type Filter interface {
	WrapService(w http.ResponseWriter, r *http.Request, s Service) error
}

type FilterFunc func(w http.ResponseWriter, r *http.Request, s Service) error

func (f FilterFunc) WrapService(w http.ResponseWriter, r *http.Request, s Service) error {
	return f(w, r, s)
}

func ComposeFilter(filters ...Filter) Filter {
	return FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		for i := len(filters) - 1; i >= 0; i-- {
			f := filters[i]
			s = ApplyFilter(f, s)
		}
		return s.ServeHTTP(w, r)
	})
}

func ApplyFilter(f Filter, s Service) Service {
	return ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		return f.WrapService(w, r, s)
	})
}

func FilterToMiddleware(f Filter, eh ErrorHandler) Middleware {
	return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		err := f.WrapService(w, r, HandlerToService(h))
		if err != nil {
			eh(w, r, err)
		}
	})
}

func HandlerToService(h http.Handler) Service {
	return ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	})
}
