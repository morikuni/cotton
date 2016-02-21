package cotton

import (
	"net/http"
)

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

type UnsupportedMethod struct {
	Method string
	Expect []string
}

func MethodFilter(methods ...string) Filter {
	ms := make(map[string]struct{})
	for _, m := range methods {
		ms[m] = struct{}{}
	}

	return func(w http.ResponseWriter, r *http.Request, s Service) Error {
		if _, ok := ms[r.Method]; !ok {
			return UnsupportedMethod{
				Method: r.Method,
				Expect: methods,
			}
		}
		return s(w, r)
	}
}
