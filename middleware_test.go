package yacm

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	assert := assert.New(t)

	count := 0
	m := Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(0, count)
			count++
			next.ServeHTTP(w, r)
		})
	})

	f := FilterFunc(func(w http.ResponseWriter, r *http.Request, s Service) error {
		assert.Equal(1, count)
		count++
		return s.TryServeHTTP(w, r)
	})

	s := ApplyFilterToFunc(ComposeFilters(m, f), func(w http.ResponseWriter, r *http.Request) error {
		assert.Equal(2, count)
		count++
		return errors.New("test")
	})

	err := s.TryServeHTTP(nil, nil)
	assert.Equal(errors.New("test"), err)
	assert.Equal(3, count)
}
