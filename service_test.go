package yacm

import (
	"net/http"
	"testing"
)

func TestService(t *testing.T) {
	count := 0
	s := ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		count++
		return nil
	})

	h := ServiceToHandler(s, func(w http.ResponseWriter, r *http.Request, err error) {
		t.Error("unreachable")
	})

	h.ServeHTTP(nil, nil)
	if count != 1 {
		t.Error("count must be 1")
	}
}
