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

func ComposeFilters(filters ...Filter) Filter {
	l := len(filters)
	switch l {
	case 0:
		panic(ErrEmptyArgs)
	case 1:
		return filters[0]
	default:
		f := filters[0]
		next := ComposeFilters(filters[1:]...)
		return FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
			return ApplyFilter(f, ApplyFilter(next, s)).TryServeHTTP(w, r)
		})
	}
}

func ApplyFilter(f Filter, s Service) Service {
	return ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		return f.WrapService(w, r, s)
	})
}
