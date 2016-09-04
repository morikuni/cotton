package yacm

import (
	"errors"
	"net/http"
	"testing"
)

func TestMiddleware_Compose(t *testing.T) {
	count := 0
	m1 := Middleware(func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		count++
		h(w, r)
	})

	m2 := m1.Compose(func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		if count != 1 {
			t.Error("count must be 1")
		}
		count++
	})

	m2(nil, nil, nil)
	if count != 2 {
		t.Error("count must be 2")
	}
}

func TestMiddleware_Apply(t *testing.T) {
	count := 0
	m := Middleware(func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		count++
		h(w, r)
	})

	h := m.Apply(func(w http.ResponseWriter, r *http.Request) {
		if count != 1 {
			t.Error("count must be 1")
		}
		count++
	})

	h(nil, nil)
	if count != 2 {
		t.Error("count must be 2")
	}
}

func TestMiddleware_ToFilter(t *testing.T) {
	count := 0
	m := Middleware(func(w http.ResponseWriter, r *http.Request, h http.HandlerFunc) {
		count++
		h(w, r)
	})

	f := m.ToFilter()

	s := f.Apply(func(w http.ResponseWriter, r *http.Request) error {
		if count != 1 {
			t.Error("count must be 1")
		}
		count++
		return errors.New("error")
	})

	err := s(nil, nil)
	if count != 2 {
		t.Error("count must be 2")
	}
	if err == nil {
		t.Error("error expected but nil")
	}
	if e := err.Error(); e != "error" {
		t.Errorf("unexpected error: %s", e)
	}
}
