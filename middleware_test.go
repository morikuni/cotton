package yacm

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestComposeMiddleware(t *testing.T) {
	assert := assert.New(t)

	count := 0
	m1 := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		assert.Equal(0, count)
		count++
		h.ServeHTTP(w, r)
	})

	m2 := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		assert.Equal(1, count)
		count++
		h.ServeHTTP(w, r)
	})

	m3 := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		assert.Equal(2, count)
		count++
	})

	ComposeMiddleware(m1, m2, m3).WrapHandler(nil, nil, nil)
	assert.Equal(3, count)
}

func TestApplyMiddleware(t *testing.T) {
	assert := assert.New(t)

	count := 0
	m := MiddlewareFunc(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		assert.Equal(0, count)
		count++
		h.ServeHTTP(w, r)
	})

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(1, count)
		count++
	})

	ApplyMiddleware(m, h).ServeHTTP(nil, nil)
	assert.Equal(2, count)
}
