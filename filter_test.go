package yacm

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFilter_Compose(t *testing.T) {
	assert := assert.New(t)

	count := 0
	f1 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(0, count)
		count++
		return s.ServeHTTP(w, r)
	})

	f2 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(1, count)
		count++
		return s.ServeHTTP(w, r)
	})

	f3 := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(2, count)
		count++
		return nil
	})

	ComposeFilter(f1, f2, f3).WrapService(nil, nil, nil)
	assert.Equal(3, count)
}

func TestFilter_Apply(t *testing.T) {
	assert := assert.New(t)

	count := 0
	f := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(0, count)
		count++
		return s.ServeHTTP(w, r)
	})

	s := ApplyFilter(f, ServiceFunc(func(w http.ResponseWriter, r *http.Request) error {
		assert.Equal(1, count)
		count++
		return nil
	}))

	err := s.ServeHTTP(nil, nil)
	assert.Equal(2, count)
	assert.Equal(nil, err)
}
