package yacm

import (
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
		return
	})

	f := m.ToFilter()

	f(nil, nil, nil)
	if count != 1 {
		t.Error("count must be 1")
	}
}
