package yacm

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestComposeErrorHanler(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	h1 := ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(0, count)
		count++
		return err
	})

	h2 := ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(dummyError, err)
		assert.Equal(1, count)
		count++
		return nil
	})

	h3 := ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Fail("unexpected call")
		return nil
	})

	err := ComposeErrorHandler(h1, h2, h3).HandleError(nil, nil, dummyError)
	assert.Equal(2, count)
	assert.Equal(nil, err)
}

func TestApplyErrorHandler(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	h := ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(0, count)
		count++
		return err
	})

	s := ErrorShutterFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		assert.Equal(dummyError, err)
		assert.Equal(1, count)
		count++
	})

	ApplyErrorHandler(h, s).ShutError(nil, nil, dummyError)
	assert.Equal(2, count)
}
