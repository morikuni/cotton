package cotton

import (
	"net/http"
)

type Service func(w http.ResponseWriter, r *http.Request) Error

func (s Service) Recover(f RecoverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s(w, r)
		if err != nil {
			f(w, r, err)
		}
	}
}

type RecoverFunc func(w http.ResponseWriter, r *http.Request, err Error)
