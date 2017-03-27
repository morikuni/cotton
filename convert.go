package yacm

import (
	"net/http"
)

func FilterToMiddleware(f Filter, eh ErrorHandler) Middleware {
	return MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		err := f.WrapService(w, r, HandlerToService(h))
		if err != nil {
			eh(w, r, err)
		}
	})
}

func MiddlewareToFilter(m Middleware) Filter {
	return FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		var err error
		m.WrapHandler(w, r, http.HandlerFunc(func(w2 http.ResponseWriter, r2 *http.Request) {
			err = s.TryServeHTTP(w2, r2)
		}))
		return err
	})
}

func ServiceToHandler(s Service, eh ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.TryServeHTTP(w, r)
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
