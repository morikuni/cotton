package yacm

import (
	"net/http"
)

// Filter is a middleware for Service.
type Filter interface {
	// WrapService handles the request and call Service.
	WrapService(w http.ResponseWriter, r *http.Request, s Service) error
}

// FilterFunc is the adapter to use a function as Filter.
type FilterFunc func(w http.ResponseWriter, r *http.Request, s Service) error

// WrapService implements Filter.
func (f FilterFunc) WrapService(w http.ResponseWriter, r *http.Request, s Service) error {
	return f(w, r, s)
}

// ComposeFilters composes multiple Filters to a single Filter.
// NOPFilter will be returned for empty arguments.
func ComposeFilters(filters ...Filter) Filter {
	l := len(filters)
	switch l {
	case 0:
		return NOPFilter
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

// ApplyFilter apply Filter to Shutter and creates a new Service.
func ApplyFilter(f Filter, s Service) Service {
	return ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		return f.WrapService(w, r, s)
	})
}

var (
	// NOPFilter is Filter that does nothing.
	NOPFilter = FilterFunc(nopFilter)
)

func nopFilter(w http.ResponseWriter, r *http.Request, s Service) error {
	return s.TryServeHTTP(w, r)
}
