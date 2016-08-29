package yacm

import (
	"net/http"
	"testing"

	"github.com/morikuni/yacm/testutil"
)

func TestFilter_Compose(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	m1 := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		callMe.Call()
		return s(w, r)
	})

	m2 := m1.Compose(func(w http.ResponseWriter, r *http.Request, s Service) error {
		callMe.MustCalled()
		callMe.Call()
		return s(w, r)
	})

	service := m2.For(testutil.NOPHandler)
	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}

func TestFilter_For(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	count := 0
	filter := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		count++
		tt.MustEqual(count, 1)
		callMe.Call()
		return s(w, r)
	})

	service := filter.For(func(w http.ResponseWriter, r *http.Request) {
		count++
		tt.MustEqual(count, 2)
		callMe.MustCalled()
		callMe.Call()
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}

func TestFilter_Then(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	count := 0
	filter := Filter(func(w http.ResponseWriter, r *http.Request, s Service) error {
		count++
		tt.MustEqual(count, 1)
		callMe.Call()
		return s(w, r)
	})

	service := filter.Then(func(w http.ResponseWriter, r *http.Request) error {
		count++
		tt.MustEqual(count, 2)
		callMe.MustCalled()
		callMe.Call()
		return nil
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}