package cotton

import (
	"net/http"
	"testing"

	"github.com/morikuni/cotton/testutil"
)

func TestMiddleware_And(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	m1 := Middleware(func(w http.ResponseWriter, r *http.Request, s Service) Error {
		callMe.Call()
		return s(w, r)
	})

	m2 := m1.And(func(w http.ResponseWriter, r *http.Request, s Service) Error {
		callMe.MustCalled()
		callMe.Call()
		return s(w, r)
	})

	service := m2.For(testutil.NOPHandler)
	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}

func TestMiddleware_For(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	count := 0
	middleware := Middleware(func(w http.ResponseWriter, r *http.Request, s Service) Error {
		count++
		tt.MustEqual(count, 1)
		callMe.Call()
		return s(w, r)
	})

	service := middleware.For(func(w http.ResponseWriter, r *http.Request) {
		count++
		tt.MustEqual(count, 2)
		callMe.MustCalled()
		callMe.Call()
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}

func TestMiddleware_Then(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	count := 0
	middleware := Middleware(func(w http.ResponseWriter, r *http.Request, s Service) Error {
		count++
		tt.MustEqual(count, 1)
		callMe.Call()
		return s(w, r)
	})

	service := middleware.Then(func(w http.ResponseWriter, r *http.Request) Error {
		count++
		tt.MustEqual(count, 2)
		callMe.MustCalled()
		callMe.Call()
		return nil
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalledTimes(2)
}
