package yacm

import (
	"net/http"
	"testing"
)

func TestFilter_Compose(t *testing.T) {
	count := 0
	f1 := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		count++
		return s(w, r)
	})

	f2 := f1.Compose(func(w http.ResponseWriter, r *http.Request, s Service) error {
		if count != 1 {
			t.Error("count must be 1")
		}
		count++
		return nil
	})

	f2(nil, nil, nil)
	if count != 2 {
		t.Error("count must be 2")
	}
}

func TestFilter_ApplyHandler(t *testing.T) {
	count := 0
	filter := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		count++
		return s(w, r)
	})

	service := filter.ApplyHandler(func(w http.ResponseWriter, r *http.Request) {
		if count != 1 {
			t.Error("count must be 1")
			count++
		}
		count++
	})

	err := service(nil, nil)
	if err != nil {
		t.Errorf("err must be nil but: %s", err.Error())
	}
	if count != 2 {
		t.Error("count must be 2")
	}
}

func TestFilter_Apply(t *testing.T) {
	count := 0
	filter := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		count++
		return s(w, r)
	})

	service := filter.Apply(func(w http.ResponseWriter, r *http.Request) error {
		if count != 1 {
			t.Error("count must be 1")
			count++
		}
		count++
		return nil
	})

	err := service(nil, nil)
	if err != nil {
		t.Errorf("err must be nil but: %s", err.Error())
	}
	if count != 2 {
		t.Error("count must be 2")
	}
}
