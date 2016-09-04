package yacm

import (
	"net/http"
	"testing"
)

func TestService(t *testing.T) {
	count := 0
	service := Service(func(w http.ResponseWriter, r *http.Request) error {
		count++
		return nil
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err error) {
		t.Error("unreachable")
	})

	handler(nil, nil)
	if count != 1 {
		t.Error("count must be 1")
	}
}
