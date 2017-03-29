package yacm

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeFilter(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	f1 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(0, count)
		count++
		return s.TryServeHTTP(w, r)
	})

	f2 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(1, count)
		count++
		return dummyError
	})

	f3 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Fail("unexpected call")
		return nil
	})

	err := ComposeFilter(f1, f2, f3).WrapService(nil, nil, nil)
	assert.Equal(2, count)
	assert.Equal(dummyError, err)
}

func TestApplyFilter(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	f := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(0, count)
		count++
		return s.TryServeHTTP(w, r)
	})

	s := ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		assert.Equal(1, count)
		count++
		return dummyError
	})

	err := ApplyFilter(f, s).TryServeHTTP(nil, nil)
	assert.Equal(2, count)
	assert.Equal(dummyError, err)
}
