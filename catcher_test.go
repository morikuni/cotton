package yacm

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeErrorHanlers(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	c1 := CatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(0, count)
		count++
		return err
	})

	c2 := CatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(dummyError, err)
		assert.Equal(1, count)
		count++
		return nil
	})

	c3 := CatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Fail("unexpected call")
		return nil
	})

	err := ComposeCatchers(c1, c2, c3).HandleError(nil, nil, dummyError)
	assert.Equal(2, count)
	assert.Equal(nil, err)
}

func TestApplyCatcher(t *testing.T) {
	assert := assert.New(t)

	dummyError := errors.New("dummy")
	count := 0
	c := CatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(0, count)
		count++
		return err
	})

	s := ErrorShutterFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		assert.Equal(dummyError, err)
		assert.Equal(1, count)
		count++
	})

	ApplyCatcher(c, s).ShutError(nil, nil, dummyError)
	assert.Equal(2, count)
}
